package service

import (
	"errors"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
)

var (
	ErrDiarioNotFound = errors.New("diário de cultivo não encontrado")
	ErrForbidden      = errors.New("acesso negado")
)

type DiarioCultivoService interface {
	CreateDiario(input dto.CreateDiarioCultivoDTO) (*dto.DiarioCultivoResponseDTO, error)
	GetDiarioByID(id uint) (*dto.DiarioCultivoResponseDTO, error)
	GetAllDiarios(page, limit int) (*dto.PaginatedResponse, error)
	UpdateDiario(id uint, input dto.UpdateDiarioCultivoDTO) (*dto.DiarioCultivoResponseDTO, error)
	DeleteDiario(id uint) error
}