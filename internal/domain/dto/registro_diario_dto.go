package dto

import (
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

// CreateRegistroDiarioDTO define a estrutura para criar um novo registro no diário.
type CreateRegistroDiarioDTO struct {
	Titulo   string              `json:"titulo" binding:"required,min=3,max=255"`
	Conteudo string              `json:"conteudo" binding:"required,min=5"`
	Data     time.Time           `json:"data" binding:"required"`
	Tipo     entity.RegistroTipo `json:"tipo" binding:"required,oneof=observacao evento aprendizado tratamento problema colheita crescimento"`
}

// RegistroDiarioResponseDTO define a estrutura para retornar um registro do diário.
type RegistroDiarioResponseDTO struct {
	ID       uint                `json:"id"`
	Titulo   string              `json:"titulo"`
	Conteudo string              `json:"conteudo"`
	Data     time.Time           `json:"data"`
	Tipo     entity.RegistroTipo `json:"tipo"`
}
