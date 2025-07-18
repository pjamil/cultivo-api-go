package service

import (
	"encoding/json"
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type GeneticaService interface {
	Criar(geneticaDto *dto.CreateGeneticaDTO) (*dto.GeneticaResponseDTO, error)
	ListarTodas(page, limit int) ([]dto.GeneticaResponseDTO, int64, error)
	BuscarPorID(id uint) (*dto.GeneticaResponseDTO, error)
	Atualizar(id uint, geneticaDto *dto.UpdateGeneticaDTO) (*dto.GeneticaResponseDTO, error)
	Deletar(id uint) error
}

// Implementação do serviço
func (s *geneticaService) Deletar(id uint) error {
	return s.repositorio.Deletar(id)
}

type geneticaService struct {
	repositorio repository.GeneticaRepositorio
}

func NewGeneticaService(repositorio repository.GeneticaRepositorio) GeneticaService {
	return &geneticaService{repositorio}
}

func (s *geneticaService) Criar(geneticaDto *dto.CreateGeneticaDTO) (*dto.GeneticaResponseDTO, error) {
	genetica := entity.Genetica{
		Nome:            geneticaDto.Nome,
		Descricao:       geneticaDto.Descricao,
		TipoGenetica:    geneticaDto.TipoGenetica,
		TipoEspecie:     geneticaDto.TipoEspecie,
		TempoFloracao:   geneticaDto.TempoFloracao,
		Origem:          geneticaDto.Origem,
		Caracteristicas: geneticaDto.Caracteristicas,
	}
	if err := s.repositorio.Criar(&genetica); err != nil {
		return nil, err
	}
	return &dto.GeneticaResponseDTO{
		ID:              genetica.ID,
		Nome:            genetica.Nome,
		Descricao:       genetica.Descricao,
		TipoGenetica:    genetica.TipoGenetica,
		TipoEspecie:     genetica.TipoEspecie,
		TempoFloracao:   genetica.TempoFloracao,
		Origem:          genetica.Origem,
		Caracteristicas: genetica.Caracteristicas,
	}, nil
}

func (s *geneticaService) ListarTodas(page, limit int) ([]dto.GeneticaResponseDTO, int64, error) {
	geneticas, total, err := s.repositorio.ListarTodos(page, limit)
	if err != nil {
		return nil, 0, err
	}

	responseDTOs := make([]dto.GeneticaResponseDTO, 0, len(geneticas))
	for _, genetica := range geneticas {
		responseDTOs = append(responseDTOs, dto.GeneticaResponseDTO{
			ID:              genetica.ID,
			Nome:            genetica.Nome,
			Descricao:       genetica.Descricao,
			TipoGenetica:    genetica.TipoGenetica,
			TipoEspecie:     genetica.TipoEspecie,
			TempoFloracao:   genetica.TempoFloracao,
			Origem:          genetica.Origem,
			Caracteristicas: genetica.Caracteristicas,
		})
	}

	return responseDTOs, total, nil
}

func (s *geneticaService) BuscarPorID(id uint) (*dto.GeneticaResponseDTO, error) {
	if id == 0 {
		return nil, gorm.ErrInvalidValue
	}

	genetica, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return nil, err
	}
	return &dto.GeneticaResponseDTO{
		ID:              genetica.ID,
		Nome:            genetica.Nome,
		Descricao:       genetica.Descricao,
		TipoGenetica:    genetica.TipoGenetica,
		TipoEspecie:     genetica.TipoEspecie,
		TempoFloracao:   genetica.TempoFloracao,
		Origem:          genetica.Origem,
		Caracteristicas: genetica.Caracteristicas,
	}, nil
}

func (s *geneticaService) Atualizar(id uint, geneticaDto *dto.UpdateGeneticaDTO) (*dto.GeneticaResponseDTO, error) {
	geneticaExistente, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	if geneticaDto.Nome != "" {
		geneticaExistente.Nome = geneticaDto.Nome
	}
	if geneticaDto.Descricao != "" {
		geneticaExistente.Descricao = geneticaDto.Descricao
	}
	if geneticaDto.TipoGenetica != "" {
		geneticaExistente.TipoGenetica = geneticaDto.TipoGenetica
	}
	if geneticaDto.TipoEspecie != "" {
		geneticaExistente.TipoEspecie = geneticaDto.TipoEspecie
	}
	if geneticaDto.TempoFloracao != 0 {
		geneticaExistente.TempoFloracao = geneticaDto.TempoFloracao
	}
	if geneticaDto.Origem != "" {
		geneticaExistente.Origem = geneticaDto.Origem
	}
	if geneticaDto.Caracteristicas != "" {
		geneticaExistente.Caracteristicas = geneticaDto.Caracteristicas
	}

	if err := s.repositorio.Atualizar(geneticaExistente); err != nil {
		return nil, err
	}

	return &dto.GeneticaResponseDTO{
		ID:              geneticaExistente.ID,
		Nome:            geneticaExistente.Nome,
		Descricao:       geneticaExistente.Descricao,
		TipoGenetica:    geneticaExistente.TipoGenetica,
		TipoEspecie:     geneticaExistente.TipoEspecie,
		TempoFloracao:   geneticaExistente.TempoFloracao,
		Origem:          geneticaExistente.Origem,
		Caracteristicas: geneticaExistente.Caracteristicas,
	}, nil
}
