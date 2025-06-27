package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type GeneticaService interface {
	Criar(genetica *models.Genetica) error
	ListarTodas() ([]models.Genetica, error)
	BuscarPorID(id uint) (*models.Genetica, error)
}

type geneticaService struct {
	repositorio repository.GeneticaRepositorio
}

func NewGeneticaService(repositorio repository.GeneticaRepositorio) GeneticaService {
	return &geneticaService{repositorio}
}

func (s *geneticaService) Criar(genetica *models.Genetica) error {
	return s.repositorio.Criar(genetica)
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
