package service

import (
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
)

type MeioCultivoService interface {
	Criar(meioCultivoDto *dto.CreateMeioCultivoDTO) (*dto.MeioCultivoResponseDTO, error)
	ListarTodos(page, limit int) ([]dto.MeioCultivoResponseDTO, int64, error)
	BuscarPorID(id uint) (*dto.MeioCultivoResponseDTO, error)
	Atualizar(id uint, meioCultivoDto *dto.UpdateMeioCultivoDTO) (*dto.MeioCultivoResponseDTO, error)
	Deletar(id uint) error
}

// Implementação do serviço
func (s *meioCultivoService) Deletar(id uint) error {
	return s.repositorio.Deletar(id)
}

type meioCultivoService struct {
	repositorio repository.MeioCultivoRepositorio
}

func NewMeioCultivoService(repositorio repository.MeioCultivoRepositorio) MeioCultivoService {
	return &meioCultivoService{repositorio}
}

func (s *meioCultivoService) Criar(meioCultivoDto *dto.CreateMeioCultivoDTO) (*dto.MeioCultivoResponseDTO, error) {
	meioCultivo := entity.MeioCultivo{
		Tipo:      meioCultivoDto.Tipo,
		Descricao: meioCultivoDto.Descricao,
	}
	if err := s.repositorio.Criar(&meioCultivo); err != nil {
		return nil, fmt.Errorf("falha ao criar meio de cultivo %s: %w", meioCultivo.Descricao, err)
	}
	return &dto.MeioCultivoResponseDTO{
		ID:        meioCultivo.ID,
		Tipo:      meioCultivo.Tipo,
		Descricao: meioCultivo.Descricao,
	}, nil
}

func (s *meioCultivoService) ListarTodos(page, limit int) ([]dto.MeioCultivoResponseDTO, int64, error) {
	meioCultivos, total, err := s.repositorio.ListarTodos(page, limit)
	if err != nil {
		return nil, 0, err
	}

	responseDTOs := make([]dto.MeioCultivoResponseDTO, 0, len(meioCultivos))
	for _, meioCultivo := range meioCultivos {
		responseDTOs = append(responseDTOs, dto.MeioCultivoResponseDTO{
			ID:        meioCultivo.ID,
			Tipo:      meioCultivo.Tipo,
			Descricao: meioCultivo.Descricao,
		})
	}

	return responseDTOs, total, nil
}

func (s *meioCultivoService) BuscarPorID(id uint) (*dto.MeioCultivoResponseDTO, error) {
	if id == 0 {
		return nil, utils.ErrInvalidInput
	}
	meioCultivo, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return nil, err
	}
	return &dto.MeioCultivoResponseDTO{
		ID:        meioCultivo.ID,
		Tipo:      meioCultivo.Tipo,
		Descricao: meioCultivo.Descricao,
	}, nil
}

func (s *meioCultivoService) Atualizar(id uint, meioCultivoDto *dto.UpdateMeioCultivoDTO) (*dto.MeioCultivoResponseDTO, error) {
	meioCultivoExistente, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar meio de cultivo com ID %d: %w", id, err)
	}

	if meioCultivoDto.Tipo != "" {
		meioCultivoExistente.Tipo = meioCultivoDto.Tipo
	}
	if meioCultivoDto.Descricao != "" {
		meioCultivoExistente.Descricao = meioCultivoDto.Descricao
	}

	if err := s.repositorio.Atualizar(meioCultivoExistente); err != nil {
		return nil, fmt.Errorf("falha ao atualizar meio de cultivo com ID %d: %w", id, err)
	}

	return &dto.MeioCultivoResponseDTO{
		ID:        meioCultivoExistente.ID,
		Tipo:      meioCultivoExistente.Tipo,
		Descricao: meioCultivoExistente.Descricao,
	}, nil
}
