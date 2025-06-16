package models

import (
	"gorm.io/gorm"
)

// Ambiente representa o ambiente onde a planta está
type Ambiente struct {
	gorm.Model
	Nome           string       `gorm:"size:100;not null" json:"nome"`
	Descricao      string       `gorm:"type:text" json:"descricao,omitempty"`
	Tipo           string       `gorm:"size:50;not null" json:"tipo"`                      // Ex: "interno", "externo", "húmido", "seco"
	Comprimento    float64      `gorm:"not null" json:"comprimento"`                       // em centímetros
	Altura         float64      `gorm:"not null" json:"altura"`                            // em centímetros
	Largura        float64      `gorm:"not null" json:"largura"`                           // em centímetros
	TempoExposicao int          `gorm:"not null" json:"tempo_exposicao"`                   // em horas
	Orientacao     string       `gorm:"size:20" json:"orientacao"`                         // norte, sul, etc.
	Fotos          []Foto       `gorm:"foreignKey:AmbienteID" json:"fotos,omitempty"`      // Fotos do ambiente
	Microclima     []Microclima `gorm:"foreignKey:AmbienteID" json:"microclima,omitempty"` // Microclimas do ambiente

	// Relacionamento com a tabela Planta
	Plantas []Planta `gorm:"foreignKey:AmbienteID"`
}
