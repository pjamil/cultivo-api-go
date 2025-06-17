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

// Planta representa uma planta no sistema
type Planta struct {
	gorm.Model
	Nome         string     `gorm:"size:255;not null" json:"nome" validate:"required,min=3,max=255"`
	ComecandoDe  string     `gorm:"size:100;not null" json:"comecando_de" validate:"required,oneof=semente clone"`      // Indica se a planta está começando de semente ou clone
	Especie      Especie    `gorm:"size:100;not null" json:"especie" validate:"required,oneof=sativa indica ruderalis"` // Especie da planta (sativa, indica, ruderalis)
	DataPlantio  *time.Time `gorm:"not null" json:"data_plantio" validate:"required"`                                   // Data do plantio
	DataColheita *time.Time `json:"data_colheita,omitempty" validate:"omitempty"`                                       // Data da colheita, se já colhida
	Status       string     `gorm:"size:100;not null" json:"status" validate:"required,oneof=ativa colhida morta"`      // Status da planta (ativa, colhida, morta)
	Notas        string     `gorm:"type:text" json:"notas,omitempty"`

	FotoCapaID    *uint       `json:"foto_capa_id,omitempty" validate:"omitempty"`                           // ID da foto de capa, se houver
	FotoCapa      *Foto       `gorm:"foreignKey:FotoCapaID" json:"foto_capa,omitempty" validate:"omitempty"` // Foto de capa da planta
	GeneticaID    uint        `gorm:"not null" json:"genetica_id" validate:"required"`                       // ID da genética da planta
	Genetica      Genetica    `gorm:"foreignKey:GeneticaID" json:"genetica" validate:"required"`             // Genética da planta
	MeioCultivoID uint        `gorm:"not null" json:"meio_cultivo_id" validate:"required"`                   // ID do meio de cultivo da planta
	MeioCultivo   MeioCultivo `gorm:"foreignKey:MeioCultivoID" json:"meio_cultivo" validate:"required"`      // Meio de cultivo da planta
	AmbienteID    uint        `gorm:"not null" json:"ambiente_id" validate:"required"`                       // ID do ambiente da planta
	Ambiente      Ambiente    `gorm:"foreignKey:AmbienteID" json:"ambiente" validate:"required"`             // Ambiente da planta
	PlantaMaeID   *uint       `json:"planta_mae_id,omitempty"`                                               // ID da planta mãe, se houver
	PlantaMae     *Planta     `gorm:"foreignKey:PlantaMaeID" json:"planta_mae,omitempty"`                    // Planta mãe, se houver

	UsuarioID uint    `json:"usuario_id"`
	Usuario   Usuario `gorm:"foreignKey:UsuarioID"`
	// Relacionamentos
	Estagios  []EstagioCrescimento `gorm:"foreignKey:PlantaID"`
	Registros []RegistroPlanta     `gorm:"foreignKey:PlantaID"`
	Tarefas   []TarefaPlanta       `gorm:"foreignKey:PlantaID"`

	Diarios []DiarioCultivo `gorm:"many2many:diario_plantas;" json:"diarios"`
}
