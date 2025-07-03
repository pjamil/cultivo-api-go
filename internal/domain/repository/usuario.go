package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type UsuarioRepositorio interface {
	Criar(usuario *models.Usuario) error
	BuscarPorID(id uint) (*models.Usuario, error)
	BuscarPorEmail(email string) (*models.Usuario, error)
	ListarTodos(page, limit int) ([]models.Usuario, int64, error)
	Atualizar(usuario *models.Usuario) error
	Deletar(id uint) error
}

type usuarioRepositorio struct {
	db *gorm.DB
}

func NewUsuarioRepositorio(db *gorm.DB) UsuarioRepositorio {
	return &usuarioRepositorio{db}
}

func (r *usuarioRepositorio) Criar(usuario *models.Usuario) error {
	return r.db.Create(usuario).Error
}

func (r *usuarioRepositorio) BuscarPorID(id uint) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.First(&usuario, id).Error
	return &usuario, err
}

func (r *usuarioRepositorio) BuscarPorEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.Where("email = ?", email).First(&usuario).Error
	return &usuario, err
}

func (r *usuarioRepositorio) ListarTodos(page, limit int) ([]models.Usuario, int64, error) {
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

func (r *usuarioRepositorio) Atualizar(usuario *models.Usuario) error {
	return r.db.Save(usuario).Error
}

func (r *usuarioRepositorio) Deletar(id uint) error {
	return r.db.Delete(&models.Usuario{}, id).Error
}