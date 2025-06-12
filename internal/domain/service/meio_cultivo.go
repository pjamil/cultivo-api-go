package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

type MeioCultivoService interface {
	CreateMeioCultivo(meioCultivo *models.MeioCultivo) error
	GetAllMeioCultivos() ([]models.MeioCultivo, error)
}

type meioCultivoService struct {
	repo repository.MeioCultivoRepository
}

func NewMeioCultivoService(repo repository.MeioCultivoRepository) MeioCultivoService {
	return &meioCultivoService{repo}
}

func (s *meioCultivoService) CreateMeioCultivo(meioCultivo *models.MeioCultivo) error {
	return s.repo.Create(meioCultivo)
}

func (s *meioCultivoService) GetAllMeioCultivos() ([]models.MeioCultivo, error) {
	var meioCultivos []models.MeioCultivo
	meioCultivos, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return meioCultivos, nil
}
