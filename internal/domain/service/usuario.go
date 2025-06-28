package service

import (
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
)

type UsuarioService interface {
	Criar(usuarioDto *dto.UsuarioCreateDTO) (*dto.UsuarioResponseDTO, error)
	BuscarPorID(id uint) (*dto.UsuarioResponseDTO, error)
	ListarTodos() ([]dto.UsuarioResponseDTO, error)
	Atualizar(id uint, usuarioDto *dto.UsuarioUpdateDTO) (*dto.UsuarioResponseDTO, error)
	Deletar(id uint) error
}

type usuarioService struct {
	repositorio repository.UsuarioRepositorio
}

func NewUsuarioService(repositorio repository.UsuarioRepositorio) UsuarioService {
	return &usuarioService{repositorio}
}

func (s *usuarioService) Criar(usuarioDto *dto.UsuarioCreateDTO) (*dto.UsuarioResponseDTO, error) {
	hash, err := utils.HashPassword(usuarioDto.Senha)
	if err != nil {
		return nil, fmt.Errorf("falha ao fazer o hash da senha do usuário %s: %w", usuarioDto.Nome, err)
	}
	usuario := models.Usuario{
		Nome:         usuarioDto.Nome,
		Email:        usuarioDto.Email,
		SenhaHash:    hash,
		Preferencias: usuarioDto.Preferencias,
	}
	if err := s.repositorio.Criar(&usuario); err != nil {
		return nil, fmt.Errorf("falha ao criar usuário %s: %w", usuarioDto.Nome, err)
	}
	return &dto.UsuarioResponseDTO{
		ID:           usuario.ID,
		Nome:         usuario.Nome,
		Email:        usuario.Email,
		Preferencias: usuario.Preferencias,
	}, nil
}

func (s *usuarioService) BuscarPorID(id uint) (*dto.UsuarioResponseDTO, error) {
	usuario, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return nil, err
	}
	return &dto.UsuarioResponseDTO{
		ID:           usuario.ID,
		Nome:         usuario.Nome,
		Email:        usuario.Email,
		Preferencias: usuario.Preferencias,
	}, nil
}

func (s *usuarioService) ListarTodos() ([]dto.UsuarioResponseDTO, error) {
	usuarios, err := s.repositorio.ListarTodos()
	if err != nil {
		return nil, err
	}

	var responseDTOs []dto.UsuarioResponseDTO
	for _, usuario := range usuarios {
		responseDTOs = append(responseDTOs, dto.UsuarioResponseDTO{
			ID:           usuario.ID,
			Nome:         usuario.Nome,
			Email:        usuario.Email,
			Preferencias: usuario.Preferencias,
		})
	}
	return responseDTOs, nil
}

// Atualizar atualiza os campos do usuário informados no DTO.
// Apenas os campos enviados (não vazios) serão atualizados.
// Se nenhum campo for enviado (todos vazios), nada será alterado no registro.
// Retorna gorm.ErrRecordNotFound se o usuário não existir.
func (s *usuarioService) Atualizar(id uint, usuarioDto *dto.UsuarioUpdateDTO) (*dto.UsuarioResponseDTO, error) {
	usuario, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar usuário com ID %d: %w", id, err)
	}
	if usuarioDto.Nome != "" {
		usuario.Nome = usuarioDto.Nome
	}
	if usuarioDto.Preferencias != "" {
		usuario.Preferencias = usuarioDto.Preferencias
	}
	if err := s.repositorio.Atualizar(usuario); err != nil {
		return nil, err
	}
	return &dto.UsuarioResponseDTO{
		ID:           usuario.ID,
		Nome:         usuario.Nome,
		Email:        usuario.Email,
		Preferencias: usuario.Preferencias,
	}, nil
}

func (s *usuarioService) Deletar(id uint) error {
	return s.repositorio.Deletar(id)
}