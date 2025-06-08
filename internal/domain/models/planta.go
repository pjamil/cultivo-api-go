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
	Nome         string     `gorm:"size:255;not null" json:"nome"`
	ComecandoDe  string     `gorm:"size:100;not null" json:"comecando_de"`
	Especie      Especie    `gorm:"size:100;not null" json:"especie"`
	DataPlantiu  *time.Time `gorm:"not null" json:"data_plantiu"`
	DataColheita *time.Time `json:"data_colheita,omitempty"`
	Status       string     `gorm:"size:100;not null" json:"status"`
	Notas        string     `gorm:"type:text" json:"notas,omitempty"`

	FotoCapaID    *uint       `json:"foto_capa_id,omitempty"` // ID da foto de capa, se houver
	FotoCapa      *Foto       `gorm:"foreignKey:FotoCapaID" json:"foto_capa,omitempty"`
	GeneticaID    uint        `json:"genetica_id"`
	Genetica      Genetica    `gorm:"foreignKey:GeneticaID"`
	MeioCultivoID uint        `json:"meio_cultivo_id"`
	MeioCultivo   MeioCultivo `gorm:"foreignKey:MeioCultivoID"`
	AmbienteID    uint        `json:"ambiente_id"`
	Ambiente      Ambiente    `gorm:"foreignKey:AmbienteID"`
	PlantaMaeID   *uint       `json:"planta_mae_id,omitempty"` // ID da planta mãe, se houver
	PlantaMae     *Planta     `gorm:"foreignKey:PlantaMaeID" json:"planta_mae,omitempty"`

	Estagios   []EstagioCrescimento `gorm:"foreignKey:PlantaID"`
	Registros  []RegistroPlanta     `gorm:"foreignKey:PlantaID"`
	Tarefas    []TarefaPlanta       `gorm:"foreignKey:PlantaID"`
	CriadoEm   time.Time            `gorm:"autoCreateTime" json:"created_at"`
	AlteradoEm time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
}
