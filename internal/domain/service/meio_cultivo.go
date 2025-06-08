package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

type MeioCultivoService interface {
	CreateMeioCultivo(meioCultivo *models.MeioCultivo) error
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
