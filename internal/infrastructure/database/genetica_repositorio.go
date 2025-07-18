package database

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gorm.io/gorm"
)

type GeneticaRepositorio struct {
	db *gorm.DB
}

func NewGeneticaRepositorio(db *gorm.DB) *GeneticaRepositorio {
	return &GeneticaRepositorio{db: db}
}

func (r *GeneticaRepositorio) Criar(genetica *entity.Genetica) error {
	return r.db.Create(genetica).Error
}

func (r *GeneticaRepositorio) BuscarPorID(id uint) (*entity.Genetica, error) {
	var genetica entity.Genetica
	err := r.db.First(&genetica, id).Error
	return &genetica, err
}

func (r *GeneticaRepositorio) ListarTodos(page, limit int) ([]entity.Genetica, int64, error) {
	var geneticas []entity.Genetica
	var total int64
	offset := (page - 1) * limit
	err := r.db.Model(&entity.Genetica{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Offset(offset).Limit(limit).Find(&geneticas).Error
	return geneticas, total, err
}

func (r *GeneticaRepositorio) Atualizar(genetica *entity.Genetica) error {
	return r.db.Save(genetica).Error
}

func (r *GeneticaRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&entity.Genetica{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
