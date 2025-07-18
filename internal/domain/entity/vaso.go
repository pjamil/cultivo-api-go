package entity

import "gorm.io/gorm"

type Vaso struct {
	gorm.Model
	Nome          string  `gorm:"size:100;not null" json:"nome"`
	Material      string  `gorm:"size:50" json:"material"` // plástico, cerâmica, etc.
	Volume        float64 `json:"volume"`                  // em litros
	Diametro      float64 `json:"diametro"`                // em cm
	Cor           string  `gorm:"size:30" json:"cor"`
	FurosDrenagem int     `json:"furos_drenagem"`
}
