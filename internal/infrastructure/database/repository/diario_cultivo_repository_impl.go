package repository

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type diarioCultivoRepositorio struct {
	db *gorm.DB
}

// NewDiarioCultivoRepositorio cria uma nova instância do repositório de DiarioCultivo.
func NewDiarioCultivoRepositorio(db *gorm.DB) repository.DiarioCultivoRepositorio {
	return &diarioCultivoRepositorio{db: db}
}

// Create insere um novo diário de cultivo no banco de dados
func (r *diarioCultivoRepositorio) Create(diarioCultivo *entity.DiarioCultivo) error {
	result := r.db.Create(diarioCultivo)
	return result.Error
}

// GetAll retorna todos os diários de cultivo com paginação
func (r *diarioCultivoRepositorio) GetAllByUserID(userID uint) ([]entity.DiarioCultivo, error) {
	var diarios []entity.DiarioCultivo
	err := r.db.Where("user_id = ?", userID).Find(&diarios).Error
	return diarios, err
}

func (r *diarioCultivoRepositorio) GetAll(page, limit int) ([]entity.DiarioCultivo, int64, error) {
	var diarios []entity.DiarioCultivo
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&entity.DiarioCultivo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Find(&diarios).Error; err != nil {
		return nil, 0, err
	}

	return diarios, total, nil
}

// Update atualiza um diário de cultivo existente
func (r *diarioCultivoRepositorio) Update(diarioCultivo *entity.DiarioCultivo) error {
	return r.db.Save(diarioCultivo).Error
}

// Delete remove um diário de cultivo pelo ID
func (r *diarioCultivoRepositorio) Delete(id uint) error {
	result := r.db.Delete(&entity.DiarioCultivo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetByID busca um diário de cultivo pelo seu ID.
func (r *diarioCultivoRepositorio) GetByID(id uint) (*entity.DiarioCultivo, error) {
	var diario entity.DiarioCultivo
	result := r.db.Preload("Plantas").Preload("Ambientes").First(&diario, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &diario, result.Error
}

func (r *diarioCultivoRepositorio) AddPlantas(diarioCultivoID uint, plantas []*entity.Planta) error {
	dc := entity.DiarioCultivo{Model: gorm.Model{ID: diarioCultivoID}}
	return r.db.Model(&dc).Association("Plantas").Append(plantas)
}

func (r *diarioCultivoRepositorio) RemovePlantas(diarioCultivoID uint, plantas []*entity.Planta) error {
	dc := entity.DiarioCultivo{Model: gorm.Model{ID: diarioCultivoID}}
	return r.db.Model(&dc).Association("Plantas").Delete(plantas)
}

func (r *diarioCultivoRepositorio) AddAmbientes(diarioCultivoID uint, ambientes []*entity.Ambiente) error {
	dc := entity.DiarioCultivo{Model: gorm.Model{ID: diarioCultivoID}}
	return r.db.Model(&dc).Association("Ambientes").Append(ambientes)
}

func (r *diarioCultivoRepositorio) RemoveAmbientes(diarioCultivoID uint, ambientes []*entity.Ambiente) error {
	dc := entity.DiarioCultivo{Model: gorm.Model{ID: diarioCultivoID}}
	return r.db.Model(&dc).Association("Ambientes").Delete(ambientes)
}