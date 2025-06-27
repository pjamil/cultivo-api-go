package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type UsuarioRepositorio interface {
	Criar(usuario *models.Usuario) error
	BuscarPorID(id uint) (*models.Usuario, error)
	ListarTodos() ([]models.Usuario, error)
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

func (r *usuarioRepositorio) ListarTodos() ([]models.Usuario, error) {
	var usuarios []models.Usuario
	err := r.db.Find(&usuarios).Error
	return usuarios, err
}

func (r *usuarioRepositorio) Atualizar(usuario *models.Usuario) error {
	return r.db.Save(usuario).Error
}

func (r *usuarioRepositorio) Deletar(id uint) error {
	return r.db.Delete(&models.Usuario{}, id).Error
}
