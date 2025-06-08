package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

	"gorm.io/gorm"
)

type MeioCultivoRepository interface {
	Create(meioCultivo *models.MeioCultivo) error
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
