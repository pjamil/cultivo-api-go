package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type GeneticaRepository interface {
	Create(genetica *models.Genetica) error
	GetAll() ([]models.Genetica, error)
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
