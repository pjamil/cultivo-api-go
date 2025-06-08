package models

import (
	"gorm.io/gorm"
)

// Ambiente representa o ambiente onde a planta está
type Ambiente struct {
	gorm.Model
	Nome      string   `gorm:"size:100;not null" json:"nome"`
	Descricao string   `gorm:"type:text" json:"descricao,omitempty"`
	Plantas   []Planta `gorm:"foreignKey:AmbienteID"`
}
