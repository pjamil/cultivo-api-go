package models

import (
	"time"

	"gorm.io/gorm"
)

// RegistroPlanta representa um registro (medição, observação, etc) da planta
type RegistroPlanta struct {
	gorm.Model
	PlantaID     uint      `json:"planta_id"`
	DataRegistro time.Time `gorm:"not null" json:"data_registro"`
	Tipo         string    `gorm:"size:100;not null" json:"tipo"`
	Valor        string    `gorm:"size:255" json:"valor,omitempty"`
	Observacao   string    `gorm:"type:text" json:"observacao,omitempty"`
}
