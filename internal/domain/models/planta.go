package models

import (
	"time"

	"gorm.io/gorm"
)

// Especie representa a especie da planta
type Especie string

const (
	EspecieSativa    Especie = "sativa"
	EspecieIndica    Especie = "indica"
	EspecieRuderalis Especie = "ruderalis"
)

type PlantaStatus string

const (
	StatusGerminating PlantaStatus = "ativa"
	StatusVegetative  PlantaStatus = "colhida"
	StatusFlowering   PlantaStatus = "morta"
)

// Planta representa uma planta no sistema
type Planta struct {
	gorm.Model
	Nome         string       `gorm:"size:255;not null" json:"nome" validate:"required,min=3,max=255"`
	ComecandoDe  string       `gorm:"size:100;not null" json:"comecando_de" validate:"required,oneof=semente clone"`
	Especie      Especie      `gorm:"size:100;not null" json:"especie" validate:"required,oneof=sativa indica ruderalis"`
	DataPlantio  *time.Time   `gorm:"not null" json:"data_plantio" validate:"required"`
	DataColheita *time.Time   `json:"data_colheita,omitempty" validate:"omitempty"`
	Status       PlantaStatus `gorm:"size:100;not null" json:"status" validate:"required,oneof=ativa colhida morta"`
	Notas        *string      `gorm:"type:text" json:"notas,omitempty"`

	FotoCapaID    *uint       `json:"foto_capa_id,omitempty"`
	FotoCapa      *Foto       `gorm:"foreignKey:FotoCapaID" json:"foto_capa,omitempty"`
	GeneticaID    uint        `gorm:"not null" json:"genetica_id"`
	Genetica      Genetica    `gorm:"foreignKey:GeneticaID" json:"genetica"`
	MeioCultivoID uint        `gorm:"not null" json:"meio_cultivo_id"`
	MeioCultivo   MeioCultivo `gorm:"foreignKey:MeioCultivoID" json:"meio_cultivo"`
	AmbienteID    uint        `gorm:"not null" json:"ambiente_id"`
	Ambiente      Ambiente    `gorm:"foreignKey:AmbienteID" json:"ambiente"`
	PlantaMaeID   *uint       `json:"planta_mae_id,omitempty"`
	PlantaMae     *Planta     `gorm:"foreignKey:PlantaMaeID" json:"planta_mae,omitempty"`

	UsuarioID uint    `gorm:"not null" json:"usuario_id"`
	Usuario   Usuario `gorm:"foreignKey:UsuarioID" json:"usuario"`

	Estagios  []EstagioCrescimento `gorm:"foreignKey:PlantaID" json:"estagios"`
	Registros []RegistroPlanta     `gorm:"foreignKey:PlantaID" json:"registros"`
	Tarefas   []TarefaPlanta       `gorm:"foreignKey:PlantaID" json:"tarefas"`

	Diarios []DiarioCultivo `gorm:"many2many:diario_plantas;" json:"diarios"`
}
