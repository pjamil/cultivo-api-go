package database

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type AmbienteRepositorio struct {
	db *gorm.DB
}

func NewAmbienteRepositorio(db *gorm.DB) *AmbienteRepositorio {
	return &AmbienteRepositorio{db: db}
}

func (r *AmbienteRepositorio) Criar(ambiente *models.Ambiente) error {
	return r.db.Create(ambiente).Error
}

func (r *AmbienteRepositorio) BuscarPorID(id uint) (*models.Ambiente, error) {
	var ambiente models.Ambiente
	err := r.db.First(&ambiente, id).Error
	return &ambiente, err
}

func (r *AmbienteRepositorio) ListarTodos(page, limit int) ([]models.Ambiente, int64, error) {
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

func (r *AmbienteRepositorio) Atualizar(ambiente *models.Ambiente) error {
	return r.db.Save(ambiente).Error
}

func (r *AmbienteRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&models.Ambiente{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
