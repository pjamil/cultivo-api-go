package entity

import (
	"time"

	"gorm.io/gorm"
)

// EstagioCrescimento representa um estágio da planta
type EstagioCrescimento struct {
	gorm.Model
	PlantaID   uint       `json:"planta_id"`
	Estagio    string     `gorm:"size:100;not null" json:"estagio"`
	DataInicio time.Time  `gorm:"not null" json:"data_inicio"`
	DataFim    *time.Time `json:"data_fim,omitempty"`
}
