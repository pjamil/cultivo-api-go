package service

import (
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type AmbienteService interface {
	Criar(ambienteDto *dto.CreateAmbienteDTO) (*models.Ambiente, error)
	ListarTodos() ([]models.Ambiente, error)
	BuscarPorID(id uint) (*models.Ambiente, error)
	Atualizar(id uint, ambienteDto *dto.UpdateAmbienteDTO) (*models.Ambiente, error)
	Deletar(id uint) error
}

// Implementação do serviço
func (s *ambienteService) Deletar(id uint) error {
	return s.repositorio.Deletar(id)
}

type ambienteService struct {
	repositorio repository.AmbienteRepositorio
}

func NewAmbienteService(repositorio repository.AmbienteRepositorio) AmbienteService {
	return &ambienteService{repositorio}
}

func (s *ambienteService) Criar(ambienteDto *dto.CreateAmbienteDTO) (*models.Ambiente, error) {
	ambiente := models.Ambiente{
		Nome:           ambienteDto.Nome,
		Descricao:      ambienteDto.Descricao,
		Tipo:           ambienteDto.Tipo,
		Comprimento:    ambienteDto.Comprimento,
		Altura:         ambienteDto.Altura,
		Largura:        ambienteDto.Largura,
		TempoExposicao: ambienteDto.TempoExposicao,
	}
	if err := s.repositorio.Criar(&ambiente); err != nil {
		return nil, fmt.Errorf("falha ao criar ambiente com nome %s: %w", ambienteDto.Nome, err)
	}
	return &ambiente, nil
}

func (s *ambienteService) ListarTodos() ([]models.Ambiente, error) {
	return s.repositorio.ListarTodos()
}

func (s *ambienteService) BuscarPorID(id uint) (*models.Ambiente, error) {
	if id == 0 {
		return nil, gorm.ErrInvalidValue
	}

	return s.repositorio.BuscarPorID(id)
}

func (s *ambienteService) Atualizar(id uint, ambienteDto *dto.UpdateAmbienteDTO) (*models.Ambiente, error) {
	ambienteExistente, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar ambiente com ID %d: %w", id, err)
	}

	if ambienteDto.Nome != "" {
		ambienteExistente.Nome = ambienteDto.Nome
	}
	if ambienteDto.Descricao != "" {
		ambienteExistente.Descricao = ambienteDto.Descricao
	}
	if ambienteDto.Tipo != "" {
		ambienteExistente.Tipo = ambienteDto.Tipo
	}
	if ambienteDto.Comprimento != 0 {
		ambienteExistente.Comprimento = ambienteDto.Comprimento
	}
	if ambienteDto.Altura != 0 {
		ambienteExistente.Altura = ambienteDto.Altura
	}
	if ambienteDto.Largura != 0 {
		ambienteExistente.Largura = ambienteDto.Largura
	}
	if ambienteDto.TempoExposicao != 0 {
		ambienteExistente.TempoExposicao = ambienteDto.TempoExposicao
	}

	if err := s.repositorio.Atualizar(ambienteExistente); err != nil {
		return nil, fmt.Errorf("falha ao atualizar ambiente com ID %d: %w", id, err)
	}

	return ambienteExistente, nil
}
