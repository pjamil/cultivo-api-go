package repository

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type GeneticaRepositorio interface {
	Criar(genetica *models.Genetica) error
	ListarTodos() ([]models.Genetica, error)
	BuscarPorID(id uint) (*models.Genetica, error)
}

type geneticaRepositorio struct {
	db *gorm.DB
}

func NewGeneticaRepositorio(db *gorm.DB) GeneticaRepositorio {
	return &geneticaRepositorio{db}
}

func (r *geneticaRepositorio) Criar(genetica *models.Genetica) error {
	return r.db.Create(genetica).Error
}

func (r *geneticaRepositorio) ListarTodos() ([]models.Genetica, error) {
	var geneticas []models.Genetica
	err := r.db.Find(&geneticas).Error
	return geneticas, err
}

func (r *geneticaRepositorio) BuscarPorID(id uint) (*models.Genetica, error) {
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
