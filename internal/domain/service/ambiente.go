package service

import (
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
)

type AmbienteService interface {
	Criar(ambienteDto *dto.CreateAmbienteDTO) (*entity.Ambiente, error)
	ListarTodos(page, limit int) ([]dto.AmbienteResponseDTO, int64, error)
	BuscarPorID(id uint) (*entity.Ambiente, error)
	Atualizar(id uint, ambienteDto *dto.UpdateAmbienteDTO) (*entity.Ambiente, error)
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

func (s *ambienteService) Criar(ambienteDto *dto.CreateAmbienteDTO) (*entity.Ambiente, error) {
	ambiente := entity.Ambiente{
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

func (s *ambienteService) ListarTodos(page, limit int) ([]dto.AmbienteResponseDTO, int64, error) {
	ambientes, total, err := s.repositorio.ListarTodos(page, limit)
	if err != nil {
		return nil, 0, err
	}

	responseDTOs := make([]dto.AmbienteResponseDTO, 0, len(ambientes))
	for _, ambiente := range ambientes {
		responseDTOs = append(responseDTOs, dto.AmbienteResponseDTO{
			ID:             ambiente.ID,
			Nome:           ambiente.Nome,
			Descricao:      ambiente.Descricao,
			Tipo:           ambiente.Tipo,
			Comprimento:    ambiente.Comprimento,
			Altura:         ambiente.Altura,
			Largura:        ambiente.Largura,
			TempoExposicao: ambiente.TempoExposicao,
			Orientacao:     ambiente.Orientacao,
		})
	}

	return responseDTOs, total, nil
}

func (s *ambienteService) BuscarPorID(id uint) (*entity.Ambiente, error) {
	if id == 0 {
		return nil, utils.ErrInvalidInput
	}

	return s.repositorio.BuscarPorID(id)
}

func (s *ambienteService) Atualizar(id uint, ambienteDto *dto.UpdateAmbienteDTO) (*entity.Ambiente, error) {
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
