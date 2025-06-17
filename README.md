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
├── go.mod                  # Dependências Go
├── Dockerfile              # Configuração Docker
└── README.md               # Documentação do projeto
