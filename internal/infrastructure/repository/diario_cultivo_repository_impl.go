package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type diarioCultivoRepository struct {
	db *gorm.DB
}

func NewDiarioCultivoRepository(db *gorm.DB) *diarioCultivoRepository {
	return &diarioCultivoRepository{db: db}
}

func (r *diarioCultivoRepository) Create(diarioCultivo *models.DiarioCultivo) error {
	return r.db.Create(diarioCultivo).Error
}

func (r *diarioCultivoRepository) GetByID(id uint) (*models.DiarioCultivo, error) {
	var diarioCultivo models.DiarioCultivo
	err := r.db.Preload("Plantas").Preload("Ambientes").First(&diarioCultivo, id).Error
	return &diarioCultivo, err
}

func (r *diarioCultivoRepository) GetAll(page, limit int) ([]models.DiarioCultivo, int64, error) {
	var diariosCultivo []models.DiarioCultivo
	var total int64

	offset := (page - 1) * limit

	r.db.Model(&models.DiarioCultivo{}).Count(&total)
	err := r.db.Preload("Plantas").Preload("Ambientes").Offset(offset).Limit(limit).Find(&diariosCultivo).Error

	return diariosCultivo, total, err
}

func (r *diarioCultivoRepository) Update(diarioCultivo *models.DiarioCultivo) error {
	return r.db.Save(diarioCultivo).Error
}

func (r *diarioCultivoRepository) Delete(id uint) error {
	return r.db.Delete(&models.DiarioCultivo{}, id).Error
}

func (r *diarioCultivoRepository) AddPlantas(diarioCultivoID uint, plantas []*models.Planta) error {
	dc := models.DiarioCultivo{Model: gorm.Model{ID: diarioCultivoID}}
	return r.db.Model(&dc).Association("Plantas").Append(plantas)
}

func (r *diarioCultivoRepository) RemovePlantas(diarioCultivoID uint, plantas []*models.Planta) error {
	dc := models.DiarioCultivo{Model: gorm.Model{ID: diarioCultivoID}}
	return r.db.Model(&dc).Association("Plantas").Delete(plantas)
}

func (r *diarioCultivoRepository) AddAmbientes(diarioCultivoID uint, ambientes []*models.Ambiente) error {
	dc := models.DiarioCultivo{Model: gorm.Model{ID: diarioCultivoID}}
	return r.db.Model(&dc).Association("Ambientes").Append(ambientes)
}

func (r *diarioCultivoRepository) RemoveAmbientes(diarioCultivoID uint, ambientes []*models.Ambiente) error {
	dc := models.DiarioCultivo{Model: gorm.Model{ID: diarioCultivoID}}
	return r.db.Model(&dc).Association("Ambientes").Delete(ambientes)
}