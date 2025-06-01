# Dockerfile e Docker Compose para API de Cultivo de Plantas em Go

Vou criar um Dockerfile multi-stage e um docker-compose.yml para sua aplicação de cultivo de plantas em Go, otimizando para produção e desenvolvimento.

## Dockerfile (Multi-stage)

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o plant-cultivation-api ./cmd/api

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/plant-cultivation-api .
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./plant-cultivation-api"]
```

## docker-compose.yml

```yaml
version: "3.8"

services:
  app:
    build:
      context: .
      target: builder # Use builder stage for development
    command: go run ./cmd/api
    volumes:
      - .:/app
      - go-mod:/go/pkg/mod
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=plant_cultivation
      - DB_PORT=5432
    depends_on:
      - db
    restart: on-failure

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=plant_cultivation
    volumes:
      - postgres-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Optional - PGAdmin for database management
  pgadmin:
    image: dpage/pgadmin4
    environment:
      - PGADMIN_DEFAULT_EMAIL=admin@plantcultivation.com
      - PGADMIN_DEFAULT_PASSWORD=admin123
    ports:
      - "5050:80"
    depends_on:
      - db

volumes:
  go-mod:
  postgres-data:
```

## docker-compose.prod.yml (Para produção)

```yaml
version: "3.8"

services:
  app:
    build:
      context: .
      # No target specified - will use final stage
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=plant_cultivation
      - DB_PORT=5432
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=plant_cultivation
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  postgres-data:
```

## Explicação da Configuração

### Dockerfile Multi-stage

1. **Estágio de Build (builder)**:

   - Usa imagem Alpine do Go para compilação
   - Instala dependências de compilação
   - Baixa os módulos Go
   - Compila a aplicação com flags de otimização

2. **Estágio de Runtime**:
   - Usa imagem Alpine mínima
   - Copia apenas o binário compilado e o .env
   - Expõe a porta 8080
   - Executa o binário

### docker-compose.yml (Desenvolvimento)

- **Serviço app**:

  - Usa estágio builder para desenvolvimento
  - Monta volume com código fonte para hot-reload
  - Configura variáveis de ambiente para desenvolvimento

- **Serviço db**:

  - PostgreSQL com volume persistente
  - Healthcheck para verificar quando está pronto

- **Serviço pgadmin** (opcional):
  - Interface web para gerenciamento do banco

### docker-compose.prod.yml (Produção)

- Otimizado para produção:
  - Sem volumes de desenvolvimento
  - Restart policies para alta disponibilidade
  - Sem ferramentas de desenvolvimento

## Como Usar

**Para desenvolvimento:**

```bash
docker-compose up -d
```

**Para produção:**

```bash
docker-compose -f docker-compose.prod.yml up -d --build
```

**Para parar:**

```bash
docker-compose down
```

**Para limpar volumes:**

```bash
docker-compose down -v
```

## Benefícios desta Configuração

1. **Multi-stage build**:

   - Imagem final pequena (~10MB)
   - Segurança reduzida (apenas o binário no runtime)
   - Dependências de compilação isoladas

2. **Ambientes separados**:

   - Configurações diferentes para dev e prod
   - Hot-reload em desenvolvimento
   - Otimizado para produção

3. **Banco de dados incluído**:
   - PostgreSQL pronto para uso
   - Persistência de dados com volumes
   - Healthcheck para inicialização ordenada

Esta configuração fornece um ambiente completo para desenvolvimento e implantação da sua API de cultivo de plantas, seguindo as melhores práticas do Docker.
