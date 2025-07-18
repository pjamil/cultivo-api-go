package service

import (
	"errors"
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

type diarioCultivoService struct {
	repo repository.DiarioCultivoRepositorio
}

func NewDiarioCultivoService(
	repo repository.DiarioCultivoRepositorio,
) *diarioCultivoService {
	return &diarioCultivoService{
		repo: repo,
	}
}

func (s *diarioCultivoService) CreateDiario(input dto.CreateDiarioCultivoDTO) (*dto.DiarioCultivoResponseDTO, error) {
	diarioCultivo := entity.DiarioCultivo{
		Nome:        input.Nome,
		DataInicio:  input.DataInicio,
		DataFim:     input.DataFim,
		UsuarioID:   input.UsuarioID,
		Privacidade: input.Privacidade,
		Tags:        input.Tags,
	}

	if err := s.repo.Create(&diarioCultivo); err != nil {
		return nil, err
	}

	// TODO: Handle PlantasIDs and AmbientesIDs for many2many relationships

	return &dto.DiarioCultivoResponseDTO{
		ID:          diarioCultivo.ID,
		Nome:        diarioCultivo.Nome,
		DataInicio:  diarioCultivo.DataInicio,
		DataFim:     diarioCultivo.DataFim,
		UsuarioID:   diarioCultivo.UsuarioID,
		Privacidade: diarioCultivo.Privacidade,
		Tags:        diarioCultivo.Tags,
		CreatedAt:   diarioCultivo.CreatedAt,
		UpdatedAt:   diarioCultivo.UpdatedAt,
	}, nil
}

func (s *diarioCultivoService) GetDiarioByID(id uint) (*dto.DiarioCultivoResponseDTO, error) {
	diarioCultivo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.DiarioCultivoResponseDTO{
		ID:          diarioCultivo.ID,
		Nome:        diarioCultivo.Nome,
		DataInicio:  diarioCultivo.DataInicio,
		DataFim:     diarioCultivo.DataFim,
		UsuarioID:   diarioCultivo.UsuarioID,
		Privacidade: diarioCultivo.Privacidade,
		Tags:        diarioCultivo.Tags,
		CreatedAt:   diarioCultivo.CreatedAt,
		UpdatedAt:   diarioCultivo.UpdatedAt,
	}, nil
}

func (s *diarioCultivoService) GetAllDiarios(page, limit int) ([]dto.DiarioCultivoResponseDTO, int64, error) {
	diariosCultivo, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	responseDTOs := make([]dto.DiarioCultivoResponseDTO, len(diariosCultivo))
	for i, diario := range diariosCultivo {
		responseDTOs[i] = dto.DiarioCultivoResponseDTO{
			ID:          diario.ID,
			Nome:        diario.Nome,
			DataInicio:  diario.DataInicio,
			DataFim:     diario.DataFim,
			UsuarioID:   diario.UsuarioID,
			Privacidade: diario.Privacidade,
			Tags:        diario.Tags,
			CreatedAt:   diario.CreatedAt,
			UpdatedAt:   diario.UpdatedAt,
		}
	}

	return responseDTOs, total, nil
}

func (s *diarioCultivoService) UpdateDiario(id uint, input dto.UpdateDiarioCultivoDTO) (*dto.DiarioCultivoResponseDTO, error) {
	
	diarioCultivo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Nome != "" {
		diarioCultivo.Nome = input.Nome
	}
	if input.DataInicio != nil {
		diarioCultivo.DataInicio = *input.DataInicio
	}
	if input.DataFim != nil {
		diarioCultivo.DataFim = input.DataFim
	}
	if input.UsuarioID != 0 {
		diarioCultivo.UsuarioID = input.UsuarioID
	}
	if input.Privacidade != "" {
		diarioCultivo.Privacidade = input.Privacidade
	}
	if input.Tags != "" {
		diarioCultivo.Tags = input.Tags
	}

	// TODO: Handle PlantasIDs and AmbientesIDs for many2many relationships

	if err := s.repo.Update(diarioCultivo); err != nil {
		return nil, err
	}

	return &dto.DiarioCultivoResponseDTO{
		ID:          diarioCultivo.ID,
		Nome:        diarioCultivo.Nome,
		DataInicio:  diarioCultivo.DataInicio,
		DataFim:     diarioCultivo.DataFim,
		UsuarioID:   diarioCultivo.UsuarioID,
		Privacidade: diarioCultivo.Privacidade,
		Tags:        diarioCultivo.Tags,
		CreatedAt:   diarioCultivo.CreatedAt,
		UpdatedAt:   diarioCultivo.UpdatedAt,
	}, nil
}

func (s *diarioCultivoService) DeleteDiario(id uint) error {
	return s.repo.Delete(id)
}