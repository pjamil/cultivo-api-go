package repository

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type AmbienteRepositorio interface {
	Criar(ambiente *models.Ambiente) error
	ListarTodos(page, limit int) ([]models.Ambiente, int64, error)
	BuscarPorID(id uint) (*models.Ambiente, error)
	Atualizar(ambiente *models.Ambiente) error
	Deletar(id uint) error
}

// Implementação do repositório
func (r *ambienteRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&models.Ambiente{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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

func (r *ambienteRepositorio) ListarTodos(page, limit int) ([]models.Ambiente, int64, error) {
	var ambientes []models.Ambiente
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Ambiente{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&ambientes).Error
	return ambientes, total, err
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

func (r *ambienteRepositorio) Atualizar(ambiente *models.Ambiente) error {
	return r.db.Save(ambiente).Error
}
