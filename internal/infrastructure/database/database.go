package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/config"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
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
	if err := db.AutoMigrate(
		&models.Ambiente{},
		&models.Anotacao{},
		&models.ClimaRegistro{},
		&models.ColecaoMidia{},
		&models.DiarioCultivo{},
		&models.EstagioCrescimento{},
		&models.Foto{},
		&models.Genetica{},
		&models.MeioCultivo{},
		&models.Microclima{},
		&models.Midia{},
		&models.Planta{},
		&models.RegistroCrescimento{},
		&models.RegistroDiario{},
		&models.RegistroPlanta{},
		&models.Substrato{},
		&models.Tarefa{},
		&models.TarefaTemplate{},
		&models.TarefaPlanta{},
		&models.Usuario{},
		&models.Vaso{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Database{DB: db}, nil
}

// Adicione no init() ou em migrations:
func SetupRelationships(db *gorm.DB) {
	// Planta pertence a Ambiente e Vaso
	db.SetupJoinTable(&models.Planta{}, "Ambiente", &models.Ambiente{})
	db.SetupJoinTable(&models.Planta{}, "Vaso", &models.Vaso{})

	// Tarefas podem ter múltiplas fotos
	db.SetupJoinTable(&models.Tarefa{}, "Fotos", &models.Foto{})

	// Índices para melhor performance
	db.Migrator().CreateIndex(&models.Tarefa{}, "status")
	db.Migrator().CreateIndex(&models.Planta{}, "status")
}
