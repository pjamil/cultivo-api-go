package models

import "gorm.io/gorm"

// Genetica representa a genética da planta
type Genetica struct {
	gorm.Model
	Nome      string   `gorm:"size:100;not null" json:"nome"`
	Descricao string   `gorm:"type:text" json:"descricao,omitempty"`
	Plantas   []Planta `gorm:"foreignKey:GeneticaID"`
}
