package repository

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type AmbienteRepository interface {
	Create(ambiente *models.Ambiente) error
	GetAll() ([]models.Ambiente, error)
	FindByID(id uint) (*models.Ambiente, error)
}

type ambienteRepository struct {
	db *gorm.DB
}

func NewAmbienteRepository(db *gorm.DB) AmbienteRepository {
	return &ambienteRepository{db}
}

func (r *ambienteRepository) Create(ambiente *models.Ambiente) error {
	return r.db.Create(ambiente).Error
}

func (r *ambienteRepository) GetAll() ([]models.Ambiente, error) {
	var ambientes []models.Ambiente
	err := r.db.Find(&ambientes).Error
	return ambientes, err
}

func (r *ambienteRepository) FindByID(id uint) (*models.Ambiente, error) {
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
