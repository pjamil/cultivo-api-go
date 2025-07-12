package database

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type UsuarioRepositorio struct {
	db *gorm.DB
}

func NewUsuarioRepositorio(db *gorm.DB) *UsuarioRepositorio {
	return &UsuarioRepositorio{db: db}
}

func (r *UsuarioRepositorio) Criar(usuario *models.Usuario) error {
	return r.db.Create(usuario).Error
}

func (r *UsuarioRepositorio) BuscarPorID(id uint) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.First(&usuario, id).Error
	return &usuario, err
}

func (r *UsuarioRepositorio) BuscarPorEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.Where("email = ?", email).First(&usuario).Error
	return &usuario, err
}

func (r *UsuarioRepositorio) ListarTodos(page, limit int) ([]models.Usuario, int64, error) {
	var usuarios []models.Usuario
	var total int64
	offset := (page - 1) * limit
	err := r.db.Model(&models.Usuario{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Offset(offset).Limit(limit).Find(&usuarios).Error
	return usuarios, total, err
}

func (r *UsuarioRepositorio) Atualizar(usuario *models.Usuario) error {
	return r.db.Save(usuario).Error
}

func (r *UsuarioRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&models.Usuario{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
