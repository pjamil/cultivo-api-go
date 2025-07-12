package database

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)


type MeioCultivoRepositorio struct {
	db *gorm.DB
}

func NewMeioCultivoRepositorio(db *gorm.DB) *MeioCultivoRepositorio {
	return &MeioCultivoRepositorio{db: db}
}

func (r *MeioCultivoRepositorio) Criar(meioCultivo *models.MeioCultivo) error {
	return r.db.Create(meioCultivo).Error
}

func (r *MeioCultivoRepositorio) BuscarPorID(id uint) (*models.MeioCultivo, error) {
	var meioCultivo models.MeioCultivo
	err := r.db.First(&meioCultivo, id).Error
	return &meioCultivo, err
}

func (r *MeioCultivoRepositorio) ListarTodos(page, limit int) ([]models.MeioCultivo, int64, error) {
	var meiosCultivo []models.MeioCultivo
	var total int64
	offset := (page - 1) * limit
	err := r.db.Model(&models.MeioCultivo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Offset(offset).Limit(limit).Find(&meiosCultivo).Error
	return meiosCultivo, total, err
}

func (r *MeioCultivoRepositorio) Atualizar(meioCultivo *models.MeioCultivo) error {
	return r.db.Save(meioCultivo).Error
}

func (r *MeioCultivoRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&models.MeioCultivo{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
