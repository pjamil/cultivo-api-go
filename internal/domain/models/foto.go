package models

import (
	"gorm.io/gorm"
)

// Foto representa uma foto no sistema
type Foto struct {
	gorm.Model
	URL string `gorm:"size:255;not null" json:"url"`
}

// PlantaFotoID representa a relação entre uma planta e suas fotos, com ID
type PlantaFotoID struct {
	gorm.Model
	PlantaID uint `json:"planta_id"`
	FotoID   uint `json:"foto_id"`
	Foto     Foto `gorm:"foreignKey:FotoID" json:"foto"`
}
