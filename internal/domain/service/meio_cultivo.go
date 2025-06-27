package service

import (
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type MeioCultivoService interface {
	Criar(meioCultivoDto *dto.CreateMeioCultivoDTO) (*models.MeioCultivo, error)
	ListarTodos() ([]models.MeioCultivo, error)
	BuscarPorID(id uint) (*models.MeioCultivo, error)
}

type meioCultivoService struct {
	repositorio repository.MeioCultivoRepositorio
}

func NewMeioCultivoService(repositorio repository.MeioCultivoRepositorio) MeioCultivoService {
	return &meioCultivoService{repositorio}
}

func (s *meioCultivoService) Criar(meioCultivoDto *dto.CreateMeioCultivoDTO) (*models.MeioCultivo, error) {
	meioCultivo := models.MeioCultivo{
		Tipo:      meioCultivoDto.Tipo,
		Descricao: meioCultivoDto.Descricao,
	}
	if err := s.repositorio.Criar(&meioCultivo); err != nil {
		return nil, fmt.Errorf("falha ao criar meio de cultivo %s: %w", meioCultivo.Descricao, err)
	}
	return &meioCultivo, nil
}

func (s *meioCultivoService) ListarTodos() ([]models.MeioCultivo, error) {
	var meioCultivos []models.MeioCultivo
	meioCultivos, err := s.repositorio.ListarTodos()
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar todos os meio de cultivo %w", err)
	}
	return meioCultivos, nil
}

func (s *meioCultivoService) BuscarPorID(id uint) (*models.MeioCultivo, error) {
	if id == 0 {
		return nil, gorm.ErrInvalidValue
	}
	return s.repositorio.BuscarPorID(id)
}
