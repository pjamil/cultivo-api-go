package repository

import (
	"errors"
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

	"gorm.io/gorm"
)

type MeioCultivoRepositorio interface {
	Criar(meioCultivo *models.MeioCultivo) error
	ListarTodos(page, limit int) ([]models.MeioCultivo, int64, error)
	BuscarPorID(id uint) (*models.MeioCultivo, error)
	Atualizar(meioCultivo *models.MeioCultivo) error
	Deletar(id uint) error
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

// ListarTodos recupera todos os registros de MeioCultivo do banco de dados com paginação.
func (r *meioCultivoRepositorio) ListarTodos(page, limit int) ([]models.MeioCultivo, int64, error) {
	var meioCultivos []models.MeioCultivo
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.MeioCultivo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&meioCultivos).Error
	if err != nil {
		return nil, 0, fmt.Errorf("falha ao buscar meio de cultivo %w", err)
	}
	return meioCultivos, total, nil
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

func (r *meioCultivoRepositorio) Atualizar(meioCultivo *models.MeioCultivo) error {
	return r.db.Save(meioCultivo).Error
}

// Deletar remove um registro de MeioCultivo do banco de dados.
func (r *meioCultivoRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&models.MeioCultivo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
