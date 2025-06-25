package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type MeioCultivoService interface {
	CreateMeioCultivo(meioCultivoDto *dto.CreateMeioCultivoDTO) (*models.MeioCultivo, error)
	GetAllMeioCultivos() ([]models.MeioCultivo, error)
	GetByID(id uint) (*models.MeioCultivo, error)
}

type meioCultivoService struct {
	repo repository.MeioCultivoRepository
}

func NewMeioCultivoService(repo repository.MeioCultivoRepository) MeioCultivoService {
	return &meioCultivoService{repo}
}

func (s *meioCultivoService) CreateMeioCultivo(meioCultivoDto *dto.CreateMeioCultivoDTO) (*models.MeioCultivo, error) {
	meioCultivo := models.MeioCultivo{
		Tipo:      meioCultivoDto.Tipo,
		Descricao: meioCultivoDto.Descricao,
	}
	if err := s.repo.Create(&meioCultivo); err != nil {
		return nil, err
	}
	return &meioCultivo, nil
}

func (s *meioCultivoService) GetAllMeioCultivos() ([]models.MeioCultivo, error) {
	var meioCultivos []models.MeioCultivo
	meioCultivos, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return meioCultivos, nil
}

func (s *meioCultivoService) GetByID(id uint) (*models.MeioCultivo, error) {
	if id == 0 {
		return nil, gorm.ErrInvalidValue
	}
	return s.repo.FindByID(id)
}
