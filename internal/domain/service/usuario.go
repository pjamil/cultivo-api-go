package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
)

type UsuarioService interface {
	Create(usuarioDto *dto.UsuarioCreateDTO) (*models.Usuario, error)
	GetByID(id uint) (*models.Usuario, error)
	GetAll() ([]models.Usuario, error)
	Update(id uint, usuarioDto *dto.UsuarioUpdateDTO) error
	Delete(id uint) error
}

type usuarioService struct {
	repo repository.UsuarioRepository
}

func NewUsuarioService(repo repository.UsuarioRepository) UsuarioService {
	return &usuarioService{repo}
}

func (s *usuarioService) Create(dto *dto.UsuarioCreateDTO) (*models.Usuario, error) {
	hash, err := utils.HashPassword(dto.Senha)
	if err != nil {
		return nil, err
	}
	usuario := models.Usuario{
		Nome:         dto.Nome,
		Email:        dto.Email,
		SenhaHash:    hash,
		Preferencias: dto.Preferencias,
	}
	if err := s.repo.Create(&usuario); err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (s *usuarioService) GetByID(id uint) (*models.Usuario, error) {
	return s.repo.FindByID(id)
}

func (s *usuarioService) GetAll() ([]models.Usuario, error) {
	return s.repo.FindAll()
}

// Update atualiza os campos do usuário informados no DTO.
// Apenas os campos enviados (não vazios) serão atualizados.
// Se nenhum campo for enviado (todos vazios), nada será alterado no registro.
// Retorna gorm.ErrRecordNotFound se o usuário não existir.
func (s *usuarioService) Update(id uint, dto *dto.UsuarioUpdateDTO) error {
	usuario, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if dto.Nome != "" {
		usuario.Nome = dto.Nome
	}
	if dto.Preferencias != "" {
		usuario.Preferencias = dto.Preferencias
	}
	return s.repo.Update(usuario)
}

func (s *usuarioService) Delete(id uint) error {
	return s.repo.Delete(id)
}
