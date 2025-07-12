package dto

import (
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
)

// CreateDiarioCultivoDTO representa os dados para criar um novo diário de cultivo
type CreateDiarioCultivoDTO struct {
	Nome        string     `json:"nome" binding:"required,min=3,max=100"`
	DataInicio  time.Time  `json:"data_inicio" binding:"required"`
	DataFim     *time.Time `json:"data_fim,omitempty"`
	Ativo       *bool      `json:"ativo,omitempty"`
	UsuarioID   uint       `json:"usuario_id" binding:"required"`
	PlantasIDs  []uint     `json:"plantas_ids,omitempty"`
	AmbientesIDs []uint     `json:"ambientes_ids,omitempty"`
	Privacidade string     `json:"privacidade,omitempty" binding:"oneof=privado publico compartilhado"`
	Tags        string     `json:"tags,omitempty" binding:"max=200"`
}

// UpdateDiarioCultivoDTO representa os dados para atualizar um diário de cultivo existente
type UpdateDiarioCultivoDTO struct {
	Nome        string     `json:"nome,omitempty" binding:"min=3,max=100"`
	DataInicio  *time.Time `json:"data_inicio,omitempty"`
	DataFim     *time.Time `json:"data_fim,omitempty"`
	Ativo       *bool      `json:"ativo,omitempty"`
	UsuarioID   uint       `json:"usuario_id,omitempty"` // Should not be updated, but kept for consistency if needed
	PlantasIDs  []uint     `json:"plantas_ids,omitempty"`
	AmbientesIDs []uint     `json:"ambientes_ids,omitempty"`
	Privacidade string     `json:"privacidade,omitempty" binding:"oneof=privado publico compartilhado"`
	Tags        string     `json:"tags,omitempty" binding:"max=200"`
}

// DiarioCultivoResponseDTO representa a resposta de um diário de cultivo
type DiarioCultivoResponseDTO struct {
	ID          uint             `json:"id"`
	Nome        string           `json:"nome"`
	DataInicio  time.Time        `json:"data_inicio"`
	DataFim     *time.Time       `json:"data_fim,omitempty"`
	Ativo       bool             `json:"ativo"`
	UsuarioID   uint             `json:"usuario_id"`
	Plantas     []models.Planta  `json:"plantas,omitempty"`
	Ambientes   []models.Ambiente `json:"ambientes,omitempty"`
	Privacidade string           `json:"privacidade"`
	Tags        string           `json:"tags"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
