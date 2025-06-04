# API REST CRUD Modular e Escalável para Cultivo de Plantas em Go

Vou implementar uma aplicação modular para gerenciamento de cultivo de plantas usando Go (Golang). A aplicação seguirá boas práticas como:

- Arquitetura limpa (Clean Architecture)
- Separação de concerns (handlers, services, repositories)
- Uso de interfaces para desacoplamento
- Configuração centralizada
- Middlewares comuns
- Documentação Swagger

## Estrutura do Projeto

```txt
/plant-cultivation-api
├── cmd
│   └── api
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── controller
│   │   ├── plant_controller.go
│   │   └── health_controller.go
│   ├── domain
│   │   ├── models
│   │   │   ├── plant.go
│   │   │   └── base.go
│   │   ├── repository
│   │   │   └── plant_repository.go
│   │   └── service
│   │       └── plant_service.go
│   ├── infrastructure
│   │   ├── database
│   │   │   └── database.go
│   │   └── server
│   │       └── server.go
│   ├── middleware
│   │   └── logging.go
│   └── utils
│       └── responses.go
├── pkg
│   └── validation
│       └── validator.go
├── go.mod
├── go.sum
└── README.md
```

## Implementação

### 1. Configuração Básica

`go.mod`:

```go
module github.com/yourusername/plant-cultivation-api

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/joho/godotenv v1.5.1
    gorm.io/gorm v1.25.5
)
```

### 2. Configuração da Aplicação

`internal/config/config.go`:

```go
package config

import (
 "log"
 "os"

 "github.com/joho/godotenv"
)

type Config struct {
 DBDriver   string
 DBHost     string
 DBPort     string
 DBUser     string
 DBPassword string
 DBName     string
 ServerPort string
}

func LoadConfig() *Config {
 if err := godotenv.Load(); err != nil {
  log.Println("No .env file found")
 }

 return &Config{
  DBDriver:   getEnv("DB_DRIVER", "postgres"),
  DBHost:     getEnv("DB_HOST", "localhost"),
  DBPort:     getEnv("DB_PORT", "5432"),
  DBUser:     getEnv("DB_USER", "postgres"),
  DBPassword: getEnv("DB_PASSWORD", "postgres"),
  DBName:     getEnv("DB_NAME", "plant_cultivation"),
  ServerPort: getEnv("SERVER_PORT", "8080"),
 }
}

func getEnv(key, defaultValue string) string {
 if value, exists := os.LookupEnv(key); exists {
  return value
 }
 return defaultValue
}
```

### 3. Modelos de Domínio

`internal/domain/models/base.go`:

```go
package models

import "time"

type Base struct {
 ID        uint      `gorm:"primaryKey" json:"id"`
 CreatedAt time.Time `json:"created_at"`
 UpdatedAt time.Time `json:"updated_at"`
}
```

`internal/domain/models/plant.go`:

```go
package models

type PlantSpecies string

const (
 SpeciesTomato    PlantSpecies = "tomato"
 SpeciesLettuce   PlantSpecies = "lettuce"
 SpeciesStrawberry PlantSpecies = "strawberry"
 SpeciesBasil     PlantSpecies = "basil"
)

type Plant struct {
 Base
 Name          string       `gorm:"size:255;not null" json:"name"`
 Species       PlantSpecies `gorm:"size:100;not null" json:"species"`
 PlantingDate  string       `gorm:"size:100;not null" json:"planting_date"`
 HarvestDate   string       `gorm:"size:100" json:"harvest_date,omitempty"`
 Status        string       `gorm:"size:100;not null" json:"status"`
 GrowthStage   string       `gorm:"size:100;not null" json:"growth_stage"`
 Notes         string       `gorm:"type:text" json:"notes,omitempty"`
}
```

### 4. Interfaces de Repositório

`internal/domain/repository/plant_repository.go`:

```go
package repository

import "github.com/yourusername/plant-cultivation-api/internal/domain/models"

type PlantRepository interface {
 Create(plant *models.Plant) error
 FindAll() ([]models.Plant, error)
 FindByID(id uint) (*models.Plant, error)
 Update(plant *models.Plant) error
 Delete(id uint) error
}
```

### 5. Serviço de Domínio

`internal/domain/service/plant_service.go`:

```go
package service

import (
 "github.com/yourusername/plant-cultivation-api/internal/domain/models"
 "github.com/yourusername/plant-cultivation-api/internal/domain/repository"
)

type PlantService struct {
 repo repository.PlantRepository
}

func NewPlantService(repo repository.PlantRepository) *PlantService {
 return &PlantService{repo: repo}
}

func (s *PlantService) CreatePlant(plant *models.Plant) error {
 return s.repo.Create(plant)
}

func (s *PlantService) GetAllPlants() ([]models.Plant, error) {
 return s.repo.FindAll()
}

func (s *PlantService) GetPlantByID(id uint) (*models.Plant, error) {
 return s.repo.FindByID(id)
}

func (s *PlantService) UpdatePlant(plant *models.Plant) error {
 return s.repo.Update(plant)
}

func (s *PlantService) DeletePlant(id uint) error {
 return s.repo.Delete(id)
}
```

### 6. Controladores

`internal/controller/plant_controller.go`:

```go
package controller

import (
 "net/http"
 "strconv"

 "github.com/gin-gonic/gin"
 "github.com/yourusername/plant-cultivation-api/internal/domain/models"
 "github.com/yourusername/plant-cultivation-api/internal/domain/service"
 "github.com/yourusername/plant-cultivation-api/internal/utils"
)

type PlantController struct {
 plantService *service.PlantService
}

func NewPlantController(plantService *service.PlantService) *PlantController {
 return &PlantController{plantService: plantService}
}

func (c *PlantController) CreatePlant(ctx *gin.Context) {
 var plant models.Plant
 if err := ctx.ShouldBindJSON(&plant); err != nil {
  utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
  return
 }

 if err := c.plantService.CreatePlant(&plant); err != nil {
  utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
  return
 }

 utils.RespondWithJSON(ctx, http.StatusCreated, plant)
}

func (c *PlantController) GetAllPlants(ctx *gin.Context) {
 plants, err := c.plantService.GetAllPlants()
 if err != nil {
  utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
  return
 }

 utils.RespondWithJSON(ctx, http.StatusOK, plants)
}

func (c *PlantController) GetPlantByID(ctx *gin.Context) {
 id, err := strconv.Atoi(ctx.Param("id"))
 if err != nil {
  utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid plant ID")
  return
 }

 plant, err := c.plantService.GetPlantByID(uint(id))
 if err != nil {
  utils.RespondWithError(ctx, http.StatusNotFound, "Plant not found")
  return
 }

 utils.RespondWithJSON(ctx, http.StatusOK, plant)
}

func (c *PlantController) UpdatePlant(ctx *gin.Context) {
 id, err := strconv.Atoi(ctx.Param("id"))
 if err != nil {
  utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid plant ID")
  return
 }

 var plant models.Plant
 if err := ctx.ShouldBindJSON(&plant); err != nil {
  utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
  return
 }

 plant.ID = uint(id)
 if err := c.plantService.UpdatePlant(&plant); err != nil {
  utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
  return
 }

 utils.RespondWithJSON(ctx, http.StatusOK, plant)
}

func (c *PlantController) DeletePlant(ctx *gin.Context) {
 id, err := strconv.Atoi(ctx.Param("id"))
 if err != nil {
  utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid plant ID")
  return
 }

 if err := c.plantService.DeletePlant(uint(id)); err != nil {
  utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
  return
 }

 utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Plant deleted successfully"})
}
```

### 7. Implementação do Repositório (usando GORM)

`internal/infrastructure/database/database.go`:

```go
package database

import (
 "fmt"
 "log"

 "gorm.io/driver/postgres"
 "gorm.io/gorm"

 "github.com/yourusername/plant-cultivation-api/internal/config"
 "github.com/yourusername/plant-cultivation-api/internal/domain/models"
)

type Database struct {
 DB *gorm.DB
}

func NewDatabase(config *config.Config) (*Database, error) {
 dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
  config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

 db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
 if err != nil {
  return nil, fmt.Errorf("failed to connect to database: %w", err)
 }

 log.Println("Database connection established")

 // Auto migrate models
 if err := db.AutoMigrate(&models.Plant{}); err != nil {
  return nil, fmt.Errorf("failed to migrate database: %w", err)
 }

 return &Database{DB: db}, nil
}

type PlantRepository struct {
 db *gorm.DB
}

func NewPlantRepository(db *gorm.DB) *PlantRepository {
 return &PlantRepository{db: db}
}

func (r *PlantRepository) Create(plant *models.Plant) error {
 return r.db.Create(plant).Error
}

func (r *PlantRepository) FindAll() ([]models.Plant, error) {
 var plants []models.Plant
 if err := r.db.Find(&plants).Error; err != nil {
  return nil, err
 }
 return plants, nil
}

func (r *PlantRepository) FindByID(id uint) (*models.Plant, error) {
 var plant models.Plant
 if err := r.db.First(&plant, id).Error; err != nil {
  return nil, err
 }
 return &plant, nil
}

func (r *PlantRepository) Update(plant *models.Plant) error {
 return r.db.Save(plant).Error
}

func (r *PlantRepository) Delete(id uint) error {
 return r.db.Delete(&models.Plant{}, id).Error
}
```

### 8. Servidor e Rotas

`internal/infrastructure/server/server.go`:

```go
package server

import (
 "github.com/yourusername/plant-cultivation-api/internal/controller"
 "github.com/yourusername/plant-cultivation-api/internal/domain/service"
 "github.com/yourusername/plant-cultivation-api/internal/infrastructure/database"
 "github.com/yourusername/plant-cultivation-api/internal/middleware"

 "github.com/gin-gonic/gin"
)

type Server struct {
 Router *gin.Engine
}

func NewServer(db *database.Database) *Server {
 router := gin.Default()

 // Middlewares
 router.Use(middleware.LoggingMiddleware())

 // Health check
 healthController := controller.NewHealthController()
 router.GET("/health", healthController.CheckHealth)

 // Plant routes
 plantRepo := database.NewPlantRepository(db.DB)
 plantService := service.NewPlantService(plantRepo)
 plantController := controller.NewPlantController(plantService)

 api := router.Group("/api/v1")
 {
  api.GET("/plants", plantController.GetAllPlants)
  api.POST("/plants", plantController.CreatePlant)
  api.GET("/plants/:id", plantController.GetPlantByID)
  api.PUT("/plants/:id", plantController.UpdatePlant)
  api.DELETE("/plants/:id", plantController.DeletePlant)
 }

 return &Server{Router: router}
}
```

### 9. Utilitários

`internal/utils/responses.go`:

```go
package utils

import "github.com/gin-gonic/gin"

func RespondWithError(c *gin.Context, code int, message string) {
 c.JSON(code, gin.H{"error": message})
}

func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
 c.JSON(code, payload)
}
```

### 10. Middlewares

`internal/middleware/logging.go`:

```go
package middleware

import (
 "log"
 "time"

 "github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
 return func(c *gin.Context) {
  start := time.Now()

  c.Next()

  duration := time.Since(start)
  log.Printf("Request - Method: %s | Path: %s | Status: %d | Duration: %v",
   c.Request.Method,
   c.Request.URL.Path,
   c.Writer.Status(),
   duration,
  )
 }
}
```

### 11. Ponto de Entrada da Aplicação

`cmd/api/main.go`:

```go
package main

import (
 "log"
 "net/http"

 "github.com/yourusername/plant-cultivation-api/internal/config"
 "github.com/yourusername/plant-cultivation-api/internal/infrastructure/database"
 "github.com/yourusername/plant-cultivation-api/internal/infrastructure/server"
)

func main() {
 // Load configuration
 cfg := config.LoadConfig()

 // Initialize database
 db, err := database.NewDatabase(cfg)
 if err != nil {
  log.Fatalf("Failed to initialize database: %v", err)
 }

 // Create server
 srv := server.NewServer(db)

 // Start server
 log.Printf("Server starting on port %s", cfg.ServerPort)
 if err := http.ListenAndServe(":"+cfg.ServerPort, srv.Router); err != nil {
  log.Fatalf("Failed to start server: %v", err)
 }
}
```

## Como Executar

1. Crie um arquivo `.env` na raiz do projeto:

```txt
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=plant_cultivation
SERVER_PORT=8080
```

2. Execute o projeto:

```bash
go run cmd/api/main.go
```

## Documentação da API

Você pode adicionar Swagger para documentação. Adicione o pacote `swaggo` e anotações aos controladores.

## Próximos Passos para Melhorias

1. **Testes**: Adicionar testes unitários e de integração
2. **Autenticação**: Implementar JWT ou OAuth2
3. **Cache**: Adicionar Redis para cache
4. **Logging**: Configurar logging estruturado (Zap ou Logrus)
5. **Monitoramento**: Integrar Prometheus e Grafana
6. **Docker**: Criar Dockerfile e docker-compose
7. **CI/CD**: Configurar pipelines GitHub Actions ou GitLab CI
8. **Documentação**: Adicionar Swagger/OpenAPI

Esta implementação fornece uma base sólida para uma API REST modular e escalável para cultivo de plantas em Go, seguindo boas práticas de desenvolvimento e arquitetura limpa.

teste polling