package dto

import "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

type CreateGeneticaDTO struct {
	Nome            string          `json:"nome" binding:"required"`
	Descricao       string          `json:"descricao"`
	TipoGenetica    string          `json:"tipoGenetica" binding:"required"`
	TipoEspecie     string          `json:"tipoEspecie" binding:"required"`
	TempoFloracao   int             `json:"tempoFloracao" binding:"required"`
	Origem          string          `json:"origem" binding:"required"`
	Caracteristicas string          `json:"caracteristicas"`
	Plantas         []models.Planta `json:"plantas,omitempty"`
}
