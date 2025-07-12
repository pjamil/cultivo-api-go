package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
)

type DiarioCultivoRepository interface {
	Create(diarioCultivo *models.DiarioCultivo) error
	GetByID(id uint) (*models.DiarioCultivo, error)
	GetAll(page, limit int) ([]models.DiarioCultivo, int64, error)
	Update(diarioCultivo *models.DiarioCultivo) error
	Delete(id uint) error
	AddPlantas(diarioCultivoID uint, plantas []*models.Planta) error
	RemovePlantas(diarioCultivoID uint, plantas []*models.Planta) error
	AddAmbientes(diarioCultivoID uint, ambientes []*models.Ambiente) error
	RemoveAmbientes(diarioCultivoID uint, ambientes []*models.Ambiente) error
}
