# 🌱 Cultivo API - Gerenciamento Modular de Plantas em Go

API RESTful modular e escalável para gerenciamento de cultivo de plantas, desenvolvida em Go (Golang) seguindo Clean Architecture e boas práticas de engenharia de software.

---

## ✨ Principais Características

- **Arquitetura Limpa (Clean Architecture)**
- **Separação de camadas:** controllers, services, repositories, models
- **Uso de interfaces para desacoplamento**
- **Configuração centralizada**
- **Middlewares reutilizáveis**
- **Documentação Swagger/OpenAPI**
- **Pronta para Docker e deploy em produção**
- **Testes e exemplos de uso via Postman/HTTP**

---

## 📁 Estrutura do Projeto

```txt
/cultivo-api-go
├── cmd/
│   └── api/                # Ponto de entrada da aplicação
├── internal/
│   ├── config/             # Configuração centralizada
│   ├── controller/         # Controllers HTTP
│   ├── domain/
│   │   ├── models/         # Modelos de domínio (Planta, Genética, Ambiente, etc)
│   │   ├── repository/     # Interfaces de repositório
│   │   └── service/        # Lógica de negócio
│   ├── infrastructure/
│   │   ├── database/       # Implementação dos repositórios (GORM)
│   │   └── server/         # Inicialização do servidor HTTP
│   ├── middleware/         # Middlewares customizados
│   └── utils/              # Utilitários (respostas, helpers)
├── pkg/                    # Pacotes externos ou utilitários
├── request/                # Exemplos de requisições HTTP (testes)
├── go.mod / go.sum         # Dependências Go
├── Dockerfile / docker-compose.yml
└── README.md
```

---

## 🚀 Como Executar Localmente

1. **Clone o repositório**

   ```bash
   git clone https://github.com/seuusuario/cultivo-api-go.git
   cd cultivo-api-go
   ```

2. **Configure o arquivo `.env`**

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=cultivo
   SERVER_PORT=8080
   ```

3. **Suba o banco de dados (opcional)**

   ```bash
   docker-compose up -d db
   ```

4. **Instale as dependências**

   ```bash
   go mod tidy
   ```

5. **Execute a aplicação**

   ```bash
   go run cmd/api/main.go
   ```

---

## 🧪 Testando a API

- Utilize o arquivo `request/api-tests.http` ou importe a [collection Postman](https://www.postman.com/) para testar os endpoints.
- Exemplos de endpoints:
  - `GET /api/v1/plants` — Lista todas as plantas
  - `POST /api/v1/plants` — Cria uma nova planta
  - `GET /api/v1/geneticas` — Lista genéticas cadastradas

---

## 🛠️ Tecnologias Utilizadas

- [Go (Golang)](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin) — Web Framework
- [GORM](https://gorm.io/) — ORM para Go
- [PostgreSQL](https://www.postgresql.org/) — Banco de dados relacional
- [Docker](https://www.docker.com/) — Containerização

---

## 📚 Documentação

- **Swagger/OpenAPI:** Em breve disponível em `/swagger/index.html`
- **Exemplos de uso:** Veja a pasta `request/` para exemplos de requisições HTTP.

---

## 📝 Próximos Passos

- [ ] Adicionar autenticação JWT
- [ ] Implementar testes automatizados
- [ ] Melhorar documentação Swagger
- [ ] Adicionar cache e monitoramento

---

## 👨‍💻 Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

---

## 📄 Licença

Este projeto está sob a licença MIT.

---

Desenvolvido por Paulo Jamil 🚀
