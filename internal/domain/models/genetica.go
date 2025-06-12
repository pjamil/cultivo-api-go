package models

import "gorm.io/gorm"

// Genetica representa a genética da planta
type Genetica struct {
	gorm.Model
	Nome            string   `gorm:"size:100;not null" json:"nome"`
	Descricao       string   `gorm:"type:text" json:"descricao,omitempty"`
	TipoGenetica    string   `gorm:"size:50;not null" json:"tipoGenetica"` // TipoGenetica representa o tipo de genética, como sativa, indica, ruderalis e hibrido.
	TipoEspecie     string   `gorm:"size:50;not null" json:"tipoEspecie"`  // TipoEspecie representa o tipo de espécie, como regular, feminizada, automatica.
	TempoFloracao   int      `gorm:"not null" json:"tempoFloracao"`        // em dias
	Origem          string   `gorm:"size:100;not null" json:"origem"`
	Caracteristicas string   `gorm:"type:text" json:"caracteristicas,omitempty"`
	Plantas         []Planta `gorm:"foreignKey:GeneticaID"`
}
