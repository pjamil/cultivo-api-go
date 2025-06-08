package models

import "gorm.io/gorm"

// MeioCultivo representa o meio de cultivo da planta
type MeioCultivo struct {
	gorm.Model
	Tipo      string   `gorm:"size:100;not null" json:"tipo"`
	Descricao string   `gorm:"type:text" json:"descricao,omitempty"`
	Plantas   []Planta `gorm:"foreignKey:MeioCultivoID"`
}
