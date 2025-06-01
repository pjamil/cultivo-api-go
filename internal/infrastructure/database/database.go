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
	if err := db.AutoMigrate(&models.Plant{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Database{DB: db}, nil
}

// Implementação do Repositório (usando GORM)
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
