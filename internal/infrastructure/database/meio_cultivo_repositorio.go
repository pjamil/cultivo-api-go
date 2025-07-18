package database

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gorm.io/gorm"
)


type MeioCultivoRepositorio struct {
	db *gorm.DB
}

func NewMeioCultivoRepositorio(db *gorm.DB) *MeioCultivoRepositorio {
	return &MeioCultivoRepositorio{db: db}
}

func (r *MeioCultivoRepositorio) Criar(meioCultivo *entity.MeioCultivo) error {
	return r.db.Create(meioCultivo).Error
}

func (r *MeioCultivoRepositorio) BuscarPorID(id uint) (*entity.MeioCultivo, error) {
	var meioCultivo entity.MeioCultivo
	err := r.db.First(&meioCultivo, id).Error
	return &meioCultivo, err
}

func (r *MeioCultivoRepositorio) ListarTodos(page, limit int) ([]entity.MeioCultivo, int64, error) {
	var meiosCultivo []entity.MeioCultivo
	var total int64
	offset := (page - 1) * limit
	err := r.db.Model(&entity.MeioCultivo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Offset(offset).Limit(limit).Find(&meiosCultivo).Error
	return meiosCultivo, total, err
}

func (r *MeioCultivoRepositorio) Atualizar(meioCultivo *entity.MeioCultivo) error {
	return r.db.Save(meioCultivo).Error
}

func (r *MeioCultivoRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&entity.MeioCultivo{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
