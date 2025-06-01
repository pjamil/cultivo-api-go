# Configurando Pipelines CI/CD para Go com GitHub Actions e GitLab CI

Vou configurar pipelines CI/CD tanto para GitHub Actions quanto para GitLab CI para sua aplicação de cultivo de plantas em Go.

## 1. GitHub Actions

Crie o arquivo `.github/workflows/go.yml`:

```yaml
name: Go CI/CD Pipeline

on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]

env:
  DOCKER_IMAGE: ghcr.io/${{ github.repository_owner }}/plant-cultivation-api
  DB_HOST: localhost
  DB_PORT: 5432
  DB_USER: postgres
  DB_PASSWORD: postgres
  DB_NAME: plant_cultivation_test

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: plant_cultivation_test
        ports:
          - 5432:5432
        options: --health-cmd pg_isready --health-interval 10s --health-timeout 5s --health-retries 5

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "1.21"

      - name: Run unit tests
        run: go test -v ./...

      - name: Run integration tests
        run: |
          go test -tags=integration -v ./...

  build-and-push:
    name: Build and Push Docker Image
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Log in to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: |
            ${{ env.DOCKER_IMAGE }}:latest
            ${{ env.DOCKER_IMAGE }}:${{ github.sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    name: Deploy to Production
    needs: build-and-push
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Install SSH key
        uses: shimataro/ssh-key-action@v2
        with:
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          known_hosts: ${{ secrets.SSH_KNOWN_HOSTS }}

      - name: Deploy to server
        run: |
          ssh -o StrictHostKeyChecking=no ${{ secrets.SSH_USER }}@${{ secrets.SSH_HOST }} << EOF
          docker pull ${{ env.DOCKER_IMAGE }}:${{ github.sha }}
          docker-compose -f docker-compose.prod.yml up -d
          EOF
```

## 2. GitLab CI

Crie o arquivo `.gitlab-ci.yml`:

```yaml
stages:
  - test
  - build
  - deploy

variables:
  DOCKER_IMAGE: registry.gitlab.com/$CI_PROJECT_PATH
  DB_HOST: postgres
  DB_PORT: 5432
  DB_USER: postgres
  DB_PASSWORD: postgres
  DB_NAME: plant_cultivation_test

services:
  - postgres:15-alpine

test:
  stage: test
  image: golang:1.21-alpine
  before_script:
    - apk add --no-cache git gcc musl-dev
    - go mod download
  script:
    - go test -v ./...
    - go test -tags=integration -v ./...

build:
  stage: build
  image: docker:24
  services:
    - docker:24-dind
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
  script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker build -t $DOCKER_IMAGE:latest .
    - docker push $DOCKER_IMAGE:latest
    - docker tag $DOCKER_IMAGE:latest $DOCKER_IMAGE:$CI_COMMIT_SHA
    - docker push $DOCKER_IMAGE:$CI_COMMIT_SHA

deploy:
  stage: deploy
  image: alpine:latest
  needs: ["build"]
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
  before_script:
    - apk add --no-cache openssh-client rsync
    - mkdir -p ~/.ssh
    - echo "$SSH_PRIVATE_KEY" > ~/.ssh/id_rsa
    - chmod 600 ~/.ssh/id_rsa
    - echo "$SSH_KNOWN_HOSTS" > ~/.ssh/known_hosts
  script:
    - ssh $SSH_USER@$SSH_HOST "
      docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY &&
      docker pull $DOCKER_IMAGE:$CI_COMMIT_SHA &&
      docker-compose -f docker-compose.prod.yml up -d
      "
```

## Configuração Necessária

### Para GitHub Actions

1. Configure secrets no repositório GitHub:

   - `SSH_PRIVATE_KEY`: Chave privada SSH para acesso ao servidor
   - `SSH_KNOWN_HOSTS`: Known hosts do servidor
   - `SSH_USER`: Usuário SSH
   - `SSH_HOST`: Host do servidor

2. Permissões:
   - Habilite write permissions para o GITHUB_TOKEN em Settings > Actions > General

### Para GitLab CI

1. Configure variables no projeto GitLab:

   - `SSH_PRIVATE_KEY`: Chave privada SSH
   - `SSH_KNOWN_HOSTS`: Known hosts
   - `SSH_USER`: Usuário SSH
   - `SSH_HOST`: Host do servidor

2. Permissões:
   - Configure o Docker-in-Docker (dind) em Settings > CI/CD > Runners

## Funcionalidades das Pipelines

1. **Testes**:

   - Executa testes unitários e de integração
   - Banco de dados PostgreSQL para testes de integração

2. **Build**:

   - Construção da imagem Docker
   - Push para o registry (GitHub Container Registry ou GitLab Registry)

3. **Deploy**:

   - SSH para o servidor de produção
   - Pull da nova imagem
   - Atualização dos containers com docker-compose

4. **Segurança**:
   - Cache de dependências
   - Credenciais protegidas por secrets
   - Execução condicional apenas no branch main

## Melhorias Possíveis

1. **Ambientes múltiplos**:

   - Adicionar stages para staging/development
   - Approvals para produção

2. **Notificações**:

   - Slack/Email sobre status do deploy

3. **Rollback automático**:

   - Verificação de health check após deploy
   - Rollback se falhar

4. **Testes adicionais**:
   - Testes de carga
   - Análise estática de código

Estas pipelines fornecem um fluxo completo de CI/CD para sua aplicação, desde os testes até o deploy em produção, seguindo as melhores práticas para projetos Go com Docker.
