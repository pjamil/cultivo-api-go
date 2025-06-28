package dto

import (
	"time"
)

// CreatePlantaDTO representa os dados para criar uma nova planta
type CreatePlantaDTO struct {
	Nome          string    `json:"nome" binding:"required"`
	ComecandoDe   string    `json:"comecando_de"`
	Especie       string    `json:"especie" binding:"required"`
	DataPlantio   time.Time `json:"data_plantio"`
	Status        string    `json:"status"`
	Notas         string    `json:"notas"`
	GeneticaID    uint      `json:"genetica_id"`
	MeioCultivoID uint      `json:"meio_cultivo_id"`
	AmbienteID    uint      `json:"ambiente_id"`
	UsuarioID     uint      `json:"usuario_id"`
}

// UpdatePlantaDTO representa os dados para atualizar uma planta existente
type UpdatePlantaDTO struct {
	Nome          string    `json:"nome"`
	ComecandoDe   string    `json:"comecando_de"`
	Especie       string    `json:"especie"`
	DataPlantio   time.Time `json:"data_plantio"`
	DataColheita  time.Time `json:"data_colheita"`
	Status        string    `json:"status"`
	EstagioCrescimento string `json:"estagio_crescimento"`
	Notas         string    `json:"notas"`
	GeneticaID    uint      `json:"genetica_id"`
	MeioCultivoID uint      `json:"meio_cultivo_id"`
	AmbienteID    uint      `json:"ambiente_id"`
	UsuarioID     uint      `json:"usuario_id"`
}
