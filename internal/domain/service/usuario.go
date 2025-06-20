package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

type UsuarioService interface {
	Create(usuario *models.Usuario) error
	GetByID(id uint) (*models.Usuario, error)
	GetAll() ([]models.Usuario, error)
	Update(usuario *models.Usuario) error
	Delete(id uint) error
}

type usuarioService struct {
	repo repository.UsuarioRepository
}

func NewUsuarioService(repo repository.UsuarioRepository) UsuarioService {
	return &usuarioService{repo}
}

func (s *usuarioService) Create(usuario *models.Usuario) error {
	return s.repo.Create(usuario)
}

func (s *usuarioService) GetByID(id uint) (*models.Usuario, error) {
	return s.repo.FindByID(id)
}

func (s *usuarioService) GetAll() ([]models.Usuario, error) {
	return s.repo.FindAll()
}

func (s *usuarioService) Update(usuario *models.Usuario) error {
	return s.repo.Update(usuario)
}

func (s *usuarioService) Delete(id uint) error {
	return s.repo.Delete(id)
}
