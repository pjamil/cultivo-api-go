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
