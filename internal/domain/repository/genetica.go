package repository

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type GeneticaRepository interface {
	Create(genetica *models.Genetica) error
	GetAll() ([]models.Genetica, error)
	FindByID(id uint) (*models.Genetica, error)
}

type geneticaRepository struct {
	db *gorm.DB
}

func NewGeneticaRepository(db *gorm.DB) GeneticaRepository {
	return &geneticaRepository{db}
}

func (r *geneticaRepository) Create(genetica *models.Genetica) error {
	return r.db.Create(genetica).Error
}

func (r *geneticaRepository) GetAll() ([]models.Genetica, error) {
	var geneticas []models.Genetica
	err := r.db.Find(&geneticas).Error
	return geneticas, err
}

func (r *geneticaRepository) FindByID(id uint) (*models.Genetica, error) {
	var genetica models.Genetica
	result := r.db.First(&genetica, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &genetica, nil
}
