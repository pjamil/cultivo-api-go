package entity

import (
	"gorm.io/gorm"
)

// Ambiente representa o ambiente onde a planta está
type Ambiente struct {
	gorm.Model
	Nome           string       `gorm:"size:100;not null" json:"nome" validate:"required"` // Nome do ambiente
	Descricao      string       `gorm:"type:text" json:"descricao,omitempty"`
	Tipo           string       `gorm:"size:50;not null" json:"tipo" validate:"required,oneof=interno externo humido seco"` // Tipo do ambiente (interno, externo, úmido, seco)
	Comprimento    float64      `gorm:"not null" json:"comprimento" validate:"required,gt=0"`                               // em centímetros
	Altura         float64      `gorm:"not null" json:"altura" validate:"required,gt=0"`                                    // em centímetros
	Largura        float64      `gorm:"not null" json:"largura" validate:"required,gt=0"`                                   // em centímetros
	TempoExposicao int          `gorm:"not null" json:"tempo_exposicao" validate:"required,gt=0"`                           // em horas
	Orientacao     string       `gorm:"size:20" json:"orientacao" validate:"required,oneof=norte sul leste oeste"`          // norte, sul, etc.
	Fotos          []Foto       `gorm:"foreignKey:AmbienteID" json:"fotos,omitempty"`                                       // Fotos do ambiente
	Microclima     []Microclima `gorm:"foreignKey:AmbienteID" json:"microclima,omitempty"`                                  // Microclimas do ambiente
	Plantas        []Planta     `gorm:"foreignKey:AmbienteID" json:"plantas,omitempty"`                                     // Plantas associadas a este ambiente
}
