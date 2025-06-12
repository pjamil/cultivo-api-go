package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

	"gorm.io/gorm"
)

type MeioCultivoRepository interface {
	Create(meioCultivo *models.MeioCultivo) error
	GetAll() ([]models.MeioCultivo, error)
}

type meioCultivoRepository struct {
	db *gorm.DB
}

func NewMeioCultivoRepository(db *gorm.DB) MeioCultivoRepository {
	return &meioCultivoRepository{db}
}

func (r *meioCultivoRepository) Create(meioCultivo *models.MeioCultivo) error {
	return r.db.Create(meioCultivo).Error
}

// GetAll retrieves all MeioCultivo records from the database.
func (r *meioCultivoRepository) GetAll() ([]models.MeioCultivo, error) {
	var meioCultivos []models.MeioCultivo
	if err := r.db.Find(&meioCultivos).Error; err != nil {
		return nil, err
	}
	return meioCultivos, nil
}
