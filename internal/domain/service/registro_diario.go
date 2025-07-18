package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
)

// RegistroDiarioService define a interface para a lógica de negócios de RegistroDiario.
type RegistroDiarioService interface {
	CriarRegistro(diarioID uint, registroDTO *dto.CreateRegistroDiarioDTO) (*dto.RegistroDiarioResponseDTO, error)
	ListarRegistrosPorDiarioID(diarioID uint, page, limit int) (*dto.PaginatedResponse, error)
}
