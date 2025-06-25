package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type GeneticaService interface {
	CreateGenetica(genetica *models.Genetica) error
	GetAll() ([]models.Genetica, error)
	GetGeneticaByID(id uint) (*models.Genetica, error)
}

type geneticaService struct {
	repo repository.GeneticaRepository
}

func NewGeneticaService(repo repository.GeneticaRepository) GeneticaService {
	return &geneticaService{repo}
}

func (s *geneticaService) CreateGenetica(genetica *models.Genetica) error {
	return s.repo.Create(genetica)
}

func (s *geneticaService) GetAll() ([]models.Genetica, error) {
	return s.repo.GetAll()
}

func (s *geneticaService) GetGeneticaByID(id uint) (*models.Genetica, error) {
	if id == 0 {
		return nil, gorm.ErrInvalidValue
	}

	return s.repo.FindByID(id)
}
