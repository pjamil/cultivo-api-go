package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gorm.io/gorm"
)

type diarioCultivoRepository struct {
	db *gorm.DB
}

func NewDiarioCultivoRepository(db *gorm.DB) *diarioCultivoRepository {
	return &diarioCultivoRepository{db: db}
}

func (r *diarioCultivoRepository) Create(diarioCultivo *entity.DiarioCultivo) error {
	return r.db.Create(diarioCultivo).Error
}

func (r *diarioCultivoRepository) GetByID(id uint) (*entity.DiarioCultivo, error) {
	var diarioCultivo entity.DiarioCultivo
	err := r.db.First(&diarioCultivo, id).Error
	return &diarioCultivo, err
}

func (r *diarioCultivoRepository) GetAll(page, limit int) ([]entity.DiarioCultivo, int64, error) {
	var diariosCultivo []entity.DiarioCultivo
	var total int64

	offset := (page - 1) * limit

	r.db.Model(&entity.DiarioCultivo{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&diariosCultivo).Error

	return diariosCultivo, total, err
}

func (r *diarioCultivoRepository) GetAllByUserID(userID uint) ([]entity.DiarioCultivo, error) {
	var diariosCultivo []entity.DiarioCultivo
	err := r.db.Where("user_id = ?", userID).Find(&diariosCultivo).Error
	return diariosCultivo, err
}

func (r *diarioCultivoRepository) Update(diarioCultivo *entity.DiarioCultivo) error {
	return r.db.Save(diarioCultivo).Error
}

func (r *diarioCultivoRepository) Delete(id uint) error {
	return r.db.Delete(&entity.DiarioCultivo{}, id).Error
}
