package dto

import (
	"time"
)

// CreatePlantaDTO representa os dados para criar uma nova planta
type CreatePlantaDTO struct {
	Nome          string    `json:"nome" binding:"required,min=3,max=100"`
	ComecandoDe   string    `json:"comecando_de" binding:"oneof='semente' 'clone' 'muda'"`
	Especie       string    `json:"especie" binding:"required,min=3,max=100"`
	DataPlantio   time.Time `json:"data_plantio" binding:"required,lte=now"`
	Status        string    `json:"status" binding:"oneof='semente' 'vegetativo' 'floracao' 'colheita' 'curando' 'finalizado' 'problema'"`
	Notas         string    `json:"notas" binding:"max=500"`
	GeneticaID    uint      `json:"genetica_id" binding:"gt=0"`
	MeioCultivoID uint      `json:"meio_cultivo_id" binding:"gt=0"`
	AmbienteID    uint      `json:"ambiente_id" binding:"gt=0"`
	UsuarioID     uint      `json:"usuario_id" binding:"gt=0"`
}

// UpdatePlantaDTO representa os dados para atualizar uma planta existente
type UpdatePlantaDTO struct {
	Nome          string    `json:"nome" binding:"min=3,max=100"`
	ComecandoDe   string    `json:"comecando_de" binding:"oneof='semente' 'clone' 'muda'"`
	Especie       string    `json:"especie" binding:"min=3,max=100"`
	DataPlantio   time.Time `json:"data_plantio" binding:"lte=now"`
	DataColheita  time.Time `json:"data_colheita" binding:"lte=now"`
	Status        string    `json:"status" binding:"oneof='semente' 'vegetativo' 'floracao' 'colheita' 'curando' 'finalizado' 'problema'"`
	EstagioCrescimento string `json:"estagio_crescimento"`
	Notas         string    `json:"notas" binding:"max=500"`
	GeneticaID    uint      `json:"genetica_id" binding:"gt=0"`
	MeioCultivoID uint      `json:"meio_cultivo_id" binding:"gt=0"`
	AmbienteID    uint      `json:"ambiente_id" binding:"gt=0"`
	UsuarioID     uint      `json:"usuario_id" binding:"gt=0"`
}

// PlantaResponseDTO representa os dados de uma planta para resposta da API
type PlantaResponseDTO struct {
	ID            uint      `json:"id"`
	Nome          string    `json:"nome"`
	ComecandoDe   string    `json:"comecando_de"`
	Especie       string    `json:"especie"`
	DataPlantio   *time.Time `json:"data_plantio"`
	DataColheita  *time.Time `json:"data_colheita,omitempty"`
	Status        string    `json:"status"`
	Notas         *string    `json:"notas,omitempty"`
	GeneticaID    uint      `json:"genetica_id"`
	MeioCultivoID uint      `json:"meio_cultivo_id"`
	AmbienteID    uint      `json:"ambiente_id"`
	UsuarioID     uint      `json:"usuario_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RegistrarFatoDTO representa os dados para registrar um fato da planta
type RegistrarFatoDTO struct {
	Tipo     string `json:"tipo" binding:"required,oneof=observacao evento aprendizado tratamento problema colheita"`
	Titulo   string `json:"titulo" binding:"required,min=3,max=100"`
	Conteudo string `json:"conteudo" binding:"required,min=5"`
}