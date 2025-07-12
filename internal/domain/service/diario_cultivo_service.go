package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
)

type DiarioCultivoService interface {
	Create(dto *dto.CreateDiarioCultivoDTO) (*dto.DiarioCultivoResponseDTO, error)
	GetByID(id uint) (*dto.DiarioCultivoResponseDTO, error)
	GetAll(page, limit int) (*dto.PaginatedResponse, error)
	Update(id uint, dto *dto.UpdateDiarioCultivoDTO) (*dto.DiarioCultivoResponseDTO, error)
	Delete(id uint) error
}
