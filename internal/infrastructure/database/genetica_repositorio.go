package database

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type GeneticaRepositorio struct {
	db *gorm.DB
}

func NewGeneticaRepositorio(db *gorm.DB) *GeneticaRepositorio {
	return &GeneticaRepositorio{db: db}
}

func (r *GeneticaRepositorio) Criar(genetica *models.Genetica) error {
	return r.db.Create(genetica).Error
}

func (r *GeneticaRepositorio) BuscarPorID(id uint) (*models.Genetica, error) {
	var genetica models.Genetica
	err := r.db.First(&genetica, id).Error
	return &genetica, err
}

func (r *GeneticaRepositorio) ListarTodos(page, limit int) ([]models.Genetica, int64, error) {
	var geneticas []models.Genetica
	var total int64
	offset := (page - 1) * limit
	err := r.db.Model(&models.Genetica{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Offset(offset).Limit(limit).Find(&geneticas).Error
	return geneticas, total, err
}

func (r *GeneticaRepositorio) Atualizar(genetica *models.Genetica) error {
	return r.db.Save(genetica).Error
}

func (r *GeneticaRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&models.Genetica{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
