package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type GeneticaService interface {
	Criar(geneticaDto *dto.CreateGeneticaDTO) (*models.Genetica, error)
	ListarTodas() ([]models.Genetica, error)
	BuscarPorID(id uint) (*models.Genetica, error)
	Atualizar(id uint, geneticaDto *dto.UpdateGeneticaDTO) (*models.Genetica, error)
}

type geneticaService struct {
	repositorio repository.GeneticaRepositorio
}

func NewGeneticaService(repositorio repository.GeneticaRepositorio) GeneticaService {
	return &geneticaService{repositorio}
}

func (s *geneticaService) Criar(geneticaDto *dto.CreateGeneticaDTO) (*models.Genetica, error) {
	genetica := models.Genetica{
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
	return &genetica, nil
}

func (s *geneticaService) ListarTodas() ([]models.Genetica, error) {
	return s.repositorio.ListarTodos()
}

func (s *geneticaService) BuscarPorID(id uint) (*models.Genetica, error) {
	if id == 0 {
		return nil, gorm.ErrInvalidValue
	}

	return s.repositorio.BuscarPorID(id)
}

func (s *geneticaService) Atualizar(id uint, geneticaDto *dto.UpdateGeneticaDTO) (*models.Genetica, error) {
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

	return geneticaExistente, nil
}
