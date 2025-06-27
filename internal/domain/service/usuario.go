package service

import (
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
)

type UsuarioService interface {
	Criar(usuarioDto *dto.UsuarioCreateDTO) (*models.Usuario, error)
	BuscarPorID(id uint) (*models.Usuario, error)
	ListarTodos() ([]models.Usuario, error)
	Atualizar(id uint, usuarioDto *dto.UsuarioUpdateDTO) error
	Deletar(id uint) error
}

type usuarioService struct {
	repositorio repository.UsuarioRepositorio
}

func NewUsuarioService(repositorio repository.UsuarioRepositorio) UsuarioService {
	return &usuarioService{repositorio}
}

func (s *usuarioService) Criar(dto *dto.UsuarioCreateDTO) (*models.Usuario, error) {
	hash, err := utils.HashPassword(dto.Senha)
	if err != nil {
		return nil, fmt.Errorf("falha ao fazer o hash da senha do usuário %s: %w", dto.Nome, err)
	}
	usuario := models.Usuario{
		Nome:         dto.Nome,
		Email:        dto.Email,
		SenhaHash:    hash,
		Preferencias: dto.Preferencias,
	}
	if err := s.repositorio.Criar(&usuario); err != nil {
		return nil, fmt.Errorf("falha ao criar usuário %s: %w", dto.Nome, err)
	}
	return &usuario, nil
}

func (s *usuarioService) BuscarPorID(id uint) (*models.Usuario, error) {
	return s.repositorio.BuscarPorID(id)
}

func (s *usuarioService) ListarTodos() ([]models.Usuario, error) {
	return s.repositorio.ListarTodos()
}

// Atualizar atualiza os campos do usuário informados no DTO.
// Apenas os campos enviados (não vazios) serão atualizados.
// Se nenhum campo for enviado (todos vazios), nada será alterado no registro.
// Retorna gorm.ErrRecordNotFound se o usuário não existir.
func (s *usuarioService) Atualizar(id uint, dto *dto.UsuarioUpdateDTO) error {
	usuario, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return fmt.Errorf("falha ao buscar usuário com ID %d: %w", id, err)
	}
	if dto.Nome != "" {
		usuario.Nome = dto.Nome
	}
	if dto.Preferencias != "" {
		usuario.Preferencias = dto.Preferencias
	}
	return s.repositorio.Atualizar(usuario)
}

func (s *usuarioService) Deletar(id uint) error {
	return s.repositorio.Deletar(id)
}
