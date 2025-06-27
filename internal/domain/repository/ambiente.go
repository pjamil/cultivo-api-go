package repository

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type AmbienteRepositorio interface {
	Criar(ambiente *models.Ambiente) error
	ListarTodos() ([]models.Ambiente, error)
	BuscarPorID(id uint) (*models.Ambiente, error)
}

type ambienteRepositorio struct {
	db *gorm.DB
}

func NewAmbienteRepositorio(db *gorm.DB) AmbienteRepositorio {
	return &ambienteRepositorio{db}
}

func (r *ambienteRepositorio) Criar(ambiente *models.Ambiente) error {
	return r.db.Create(ambiente).Error
}

func (r *ambienteRepositorio) ListarTodos() ([]models.Ambiente, error) {
	var ambientes []models.Ambiente
	err := r.db.Find(&ambientes).Error
	return ambientes, err
}

func (r *ambienteRepositorio) BuscarPorID(id uint) (*models.Ambiente, error) {
	var ambiente models.Ambiente
	result := r.db.First(&ambiente, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &ambiente, nil
}
