package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type UsuarioRepository interface {
	Create(usuario *models.Usuario) error
	FindByID(id uint) (*models.Usuario, error)
	FindAll() ([]models.Usuario, error)
	Update(usuario *models.Usuario) error
	Delete(id uint) error
}

type usuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) UsuarioRepository {
	return &usuarioRepository{db}
}

func (r *usuarioRepository) Create(usuario *models.Usuario) error {
	return r.db.Create(usuario).Error
}

func (r *usuarioRepository) FindByID(id uint) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.First(&usuario, id).Error
	return &usuario, err
}

func (r *usuarioRepository) FindAll() ([]models.Usuario, error) {
	var usuarios []models.Usuario
	err := r.db.Find(&usuarios).Error
	return usuarios, err
}

func (r *usuarioRepository) Update(usuario *models.Usuario) error {
	return r.db.Save(usuario).Error
}

func (r *usuarioRepository) Delete(id uint) error {
	return r.db.Delete(&models.Usuario{}, id).Error
}
