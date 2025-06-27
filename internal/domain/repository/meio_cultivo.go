package repository

import (
	"errors"
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

	"gorm.io/gorm"
)

type MeioCultivoRepositorio interface {
	Criar(meioCultivo *models.MeioCultivo) error
	ListarTodos() ([]models.MeioCultivo, error)
	BuscarPorID(id uint) (*models.MeioCultivo, error)
}

type meioCultivoRepositorio struct {
	db *gorm.DB
}

func NewMeioCultivoRepositorio(db *gorm.DB) MeioCultivoRepositorio {
	return &meioCultivoRepositorio{db}
}

func (r *meioCultivoRepositorio) Criar(meioCultivo *models.MeioCultivo) error {
	return r.db.Create(meioCultivo).Error
}

// ListarTodos recupera todos os registros de MeioCultivo do banco de dados.
func (r *meioCultivoRepositorio) ListarTodos() ([]models.MeioCultivo, error) {
	var meioCultivos []models.MeioCultivo
	if err := r.db.Find(&meioCultivos).Error; err != nil {
		return nil, fmt.Errorf("falha ao buscar meio de cultivo %w", err)
	}
	return meioCultivos, nil
}

// BuscarPorID recupera um registro de MeioCultivo por seu ID.
func (r *meioCultivoRepositorio) BuscarPorID(id uint) (*models.MeioCultivo, error) {
	var meioCultivo models.MeioCultivo
	result := r.db.First(&meioCultivo, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &meioCultivo, nil
}
