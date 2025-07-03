package repository

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type GeneticaRepositorio interface {
	Criar(genetica *models.Genetica) error
	ListarTodos(page, limit int) ([]models.Genetica, int64, error)
	BuscarPorID(id uint) (*models.Genetica, error)
	Atualizar(genetica *models.Genetica) error
	Deletar(id uint) error
}

// Implementação do repositório
func (r *geneticaRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&models.Genetica{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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

func (r *geneticaRepositorio) ListarTodos(page, limit int) ([]models.Genetica, int64, error) {
	var geneticas []models.Genetica
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Genetica{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&geneticas).Error
	return geneticas, total, err
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

func (r *geneticaRepositorio) Atualizar(genetica *models.Genetica) error {
	return r.db.Save(genetica).Error
}
