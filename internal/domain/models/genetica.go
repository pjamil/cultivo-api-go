package models

import "gorm.io/gorm"

// Genetica representa a genética da planta
type Genetica struct {
	gorm.Model
	Nome            string   `gorm:"size:100;not null" json:"nome" validate:"required, min=2, max=100"` // Nome da genética
	Descricao       string   `gorm:"type:text" json:"descricao,omitempty"`
	TipoGenetica    string   `gorm:"size:50;not null" json:"tipoGenetica" validate:"required,oneof=sativa indica ruderalis hibrido"` // TipoGenetica representa o tipo de genética, como sativa, indica, ruderalis e hibrido.
	TipoEspecie     string   `gorm:"size:50;not null" json:"tipoEspecie" validate:"required,oneof=regular feminizada automatica"`    // TipoEspecie representa o tipo de espécie, como regular, feminizada, automatica.
	TempoFloracao   int      `gorm:"not null" json:"tempoFloracao" validate:"required,gt=0"`                                         // em dias
	Origem          string   `gorm:"size:100;not null" json:"origem" validate:"required, min=2, max=100"`                            // Origem da genética
	Caracteristicas string   `gorm:"type:text" json:"caracteristicas,omitempty"`
	Plantas         []Planta `gorm:"foreignKey:GeneticaID"`
}
