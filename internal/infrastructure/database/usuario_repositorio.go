package database

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gorm.io/gorm"
)

type UsuarioRepositorio struct {
	db *gorm.DB
}

func NewUsuarioRepositorio(db *gorm.DB) *UsuarioRepositorio {
	return &UsuarioRepositorio{db: db}
}

func (r *UsuarioRepositorio) Criar(usuario *entity.Usuario) error {
	return r.db.Create(usuario).Error
}

func (r *UsuarioRepositorio) BuscarPorID(id uint) (*entity.Usuario, error) {
	var usuario entity.Usuario
	err := r.db.First(&usuario, id).Error
	return &usuario, err
}

func (r *UsuarioRepositorio) BuscarPorEmail(email string) (*entity.Usuario, error) {
	var usuario entity.Usuario
	err := r.db.Where("email = ?", email).First(&usuario).Error
	return &usuario, err
}

func (r *UsuarioRepositorio) ListarTodos(page, limit int) ([]entity.Usuario, int64, error) {
	var usuarios []entity.Usuario
	var total int64
	offset := (page - 1) * limit
	err := r.db.Model(&entity.Usuario{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Offset(offset).Limit(limit).Find(&usuarios).Error
	return usuarios, total, err
}

func (r *UsuarioRepositorio) Atualizar(usuario *entity.Usuario) error {
	return r.db.Save(usuario).Error
}

func (r *UsuarioRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&entity.Usuario{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *UsuarioRepositorio) ExistePorEmail(email string) bool {
	var count int64
	r.db.Model(&entity.Usuario{}).Where("email = ?", email).Count(&count)
	return count > 0
}
