# 🌿 Cultivo API Go

[![Go Report Card](https://goreportcard.com/badge/gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go)](https://goreportcard.com/report/gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## 📝 Descrição

A `cultivo-api-go` é uma API robusta e escalável desenvolvida em Go para gerenciar dados relacionados ao cultivo de plantas. Ela oferece funcionalidades para o controle de ambientes, plantas, genéticas, usuários e muito mais, facilitando o registro e acompanhamento de todo o ciclo de cultivo.

## ✨ Funcionalidades

*   **Gerenciamento de Ambientes:** Crie, visualize, atualize e exclua informações sobre os ambientes de cultivo (interno, externo, úmido, seco).
*   **Controle de Plantas:** Registre e acompanhe suas plantas desde a germinação até a colheita, incluindo detalhes como espécie, genética, meio de cultivo e status.
*   **Dados de Genética:** Mantenha um registro detalhado das genéticas das suas plantas, incluindo tipo, origem e características.
*   **Gestão de Usuários:** Gerencie usuários com autenticação segura e preferências personalizáveis.
*   **Registro de Microclima:** Monitore e registre dados de microclima (temperatura, umidade, luminosidade) para cada ambiente.
*   **Upload de Fotos:** Associe fotos a ambientes e plantas para um registro visual completo.
*   **Migrações de Banco de Dados:** Gerenciamento de esquema de banco de dados versionado com `golang-migrate/migrate`.

## 🚀 Tecnologias Utilizadas

*   **Go (Golang):** Linguagem de programação principal.
*   **Gin Gonic:** Framework web para construção da API RESTful.
*   **GORM:** ORM (Object-Relational Mapping) para interação com o banco de dados.
*   **PostgreSQL:** Banco de dados relacional.
*   **golang-migrate/migrate:** Ferramenta para gerenciamento de migrações de banco de dados.
*   **Logrus:** Biblioteca de logging estruturado.
*   **Swagger:** Geração automática de documentação da API.

## ⚙️ Como Começar

Siga estas instruções para configurar e executar o projeto em seu ambiente local.

### Pré-requisitos

Certifique-se de ter o seguinte instalado em sua máquina:

*   [Go](https://golang.org/doc/install) (versão 1.24 ou superior)
*   [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/)
*   [make](https://www.gnu.org/software/make/manual/make.html) (geralmente já vem com sistemas Unix/Linux)

### Clonando o Repositório

```bash
git clone https://gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go.git
cd cultivo-api-go
```

### Executando com Docker Compose (Recomendado)

Esta é a maneira mais fácil de colocar a API e o banco de dados em funcionamento.

1.  **Construa e inicie os serviços:**
    ```bash
    docker-compose up --build
    ```
    Isso construirá as imagens Docker, iniciará o banco de dados PostgreSQL e a API. As migrações do banco de dados serão aplicadas automaticamente na inicialização do serviço `app`.

2.  **Acesse a API:**
    A API estará disponível em `http://localhost:8080`.
    A documentação Swagger estará em `http://localhost:8080/swagger/index.html`.

### Executando Localmente (sem Docker para a API)

Se você preferir executar a API Go diretamente em sua máquina, mas ainda usar o banco de dados via Docker Compose:

1.  **Inicie apenas o serviço de banco de dados:**
    ```bash
    docker-compose up -d db
    ```

2.  **Crie o banco de dados (se ainda não existir) e aplique as migrações:**
    Se você precisar recriar o banco de dados do zero (por exemplo, após um erro de migração ou para um ambiente limpo):
    *   Pare o serviço `app` (se estiver rodando):
        ```bash
        docker-compose stop app
        ```
    *   Acesse o shell do contêiner do banco de dados:
        ```bash
        docker-compose exec db psql -U postgres
        ```
    *   Dentro do shell do PostgreSQL, remova e recrie o banco de dados `cultivo-api-go`:
        ```sql
        DROP DATABASE IF EXISTS "cultivo-api-go";
        CREATE DATABASE "cultivo-api-go";
        \q
        ```
    *   Aplique as migrações usando o `Makefile`:
        ```bash
        make migrate-up
        ```

3.  **Execute a aplicação Go:**
    ```bash
    go run cmd/cultivo-api-go/main.go
    ```

## 🔄 Migrações de Banco de Dados

Para detalhes sobre como criar, aplicar e reverter migrações de banco de dados usando `golang-migrate/migrate`, consulte o arquivo [MIGRATIONS.md](MIGRATIONS.md).

## 📄 Documentação da API

A documentação interativa da API é gerada automaticamente usando Swagger e pode ser acessada em `http://localhost:8080/swagger/index.html` quando a aplicação estiver em execução.

## 🧪 Testes

Para executar os testes do projeto, utilize o comando:

```bash
make test
```

## 🤝 Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues e pull requests.

## 📄 Licença

Este projeto está licenciado sob a Licença Apache 2.0. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.