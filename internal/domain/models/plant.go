package models

import "gorm.io/gorm"

// PlantSpecies represents the species of a plant
// @Schema
type PlantSpecies string

const (
	// @Enum tomato, lettuce, strawberry, basil
	SpeciesTomato     PlantSpecies = "tomato"
	SpeciesLettuce    PlantSpecies = "lettuce"
	SpeciesStrawberry PlantSpecies = "strawberry"
	SpeciesBasil      PlantSpecies = "basil"
)

// Plant represents a plant in the cultivation system
// @Schema
type Plant struct {
	gorm.Model
	// @Example Tomato Plant
	Name string `gorm:"size:255;not null" json:"name"`
	// @Example tomato
	Species PlantSpecies `gorm:"size:100;not null" json:"species"`
	// @Example 2023-03-01
	PlantingDate string `gorm:"size:100;not null" json:"planting_date"`
	// @Example 2023-06-01
	HarvestDate string `gorm:"size:100" json:"harvest_date,omitempty"`
	// @Enum growing,harvested,wilted
	// @Example growing
	Status string `gorm:"size:100;not null" json:"status"`
	// @Enum seedling,vegetative,flowering,fruiting,mature
	// @Example vegetative
	GrowthStage string `gorm:"size:100;not null" json:"growth_stage"`
	// @Example Healthy plant with no signs of disease
	Notes string `gorm:"type:text" json:"notes,omitempty"`
}
