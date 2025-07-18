package dto

import (
	"time"
)

// CreateDiarioCultivoDTO representa os dados para criar um novo diário de cultivo
type CreateDiarioCultivoDTO struct {
	Nome        string    `json:"nome" binding:"required,min=3,max=255"`
	DataInicio  time.Time `json:"data_inicio" binding:"required"`
	DataFim     *time.Time `json:"data_fim"`
	UsuarioID   uint      `json:"usuario_id" binding:"required"`
	PlantasIDs  []uint    `json:"plantas_ids"`
	AmbientesIDs []uint    `json:"ambientes_ids"`
	Privacidade string    `json:"privacidade" binding:"oneof=publico privado"`
	Tags        string    `json:"tags"`
}

// UpdateDiarioCultivoDTO representa os dados para atualizar um diário de cultivo existente
type UpdateDiarioCultivoDTO struct {
	Nome        string     `json:"nome,omitempty" binding:"min=3,max=255"`
	DataInicio  *time.Time `json:"data_inicio,omitempty"`
	DataFim     *time.Time `json:"data_fim,omitempty"`
	UsuarioID   uint       `json:"usuario_id,omitempty"`
	PlantasIDs  []uint     `json:"plantas_ids,omitempty"`
	AmbientesIDs []uint     `json:"ambientes_ids,omitempty"`
	Privacidade string     `json:"privacidade,omitempty" binding:"oneof=publico privado"`
	Tags        string     `json:"tags,omitempty"`
}

// DiarioCultivoResponseDTO representa a resposta de um diário de cultivo
type DiarioCultivoResponseDTO struct {
	ID          uint      `json:"id"`
	Nome        string    `json:"nome"`
	DataInicio  time.Time `json:"data_inicio"`
	DataFim     *time.Time `json:"data_fim"`
	UsuarioID   uint      `json:"usuario_id"`
	Plantas     []PlantaResponseDTO `json:"plantas"`
	Ambientes   []AmbienteResponseDTO `json:"ambientes"`
	Privacidade string    `json:"privacidade"`
	Tags        string    `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
