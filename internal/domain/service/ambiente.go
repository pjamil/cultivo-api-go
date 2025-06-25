package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type AmbienteService interface {
	CreateAmbiente(ambienteDto *dto.CreateAmbienteDTO) (*models.Ambiente, error)
	GetAll() ([]models.Ambiente, error)
	GetAmbienteByID(id uint) (*models.Ambiente, error)
}

type ambienteService struct {
	repo repository.AmbienteRepository
}

func NewAmbienteService(repo repository.AmbienteRepository) AmbienteService {
	return &ambienteService{repo}
}

func (s *ambienteService) CreateAmbiente(ambienteDto *dto.CreateAmbienteDTO) (*models.Ambiente, error) {
	ambiente := models.Ambiente{
		Nome:           ambienteDto.Nome,
		Descricao:      ambienteDto.Descricao,
		Tipo:           ambienteDto.Tipo,
		Comprimento:    ambienteDto.Comprimento,
		Altura:         ambienteDto.Altura,
		Largura:        ambienteDto.Largura,
		TempoExposicao: ambienteDto.TempoExposicao,
	}
	if err := s.repo.Create(&ambiente); err != nil {
		return nil, err
	}
	return &ambiente, nil
}

func (s *ambienteService) GetAll() ([]models.Ambiente, error) {
	return s.repo.GetAll()
}

func (s *ambienteService) GetAmbienteByID(id uint) (*models.Ambiente, error) {
	if id == 0 {
		return nil, gorm.ErrInvalidValue
	}

	return s.repo.FindByID(id)
}
