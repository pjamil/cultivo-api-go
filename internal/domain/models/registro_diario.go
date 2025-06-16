package models

import (
	"time"

	"gorm.io/gorm"
)

type RegistroDiario struct {
	gorm.Model
	Data time.Time `gorm:"index" json:"data"`
	Tipo string    `gorm:"size:50;index" json:"tipo"` // observacao, evento, aprendizado

	DiarioCultivoID uint `json:"diario_cultivo_id"`

	// Campos polimórficos
	Titulo         string  `gorm:"size:100" json:"titulo"`
	Conteudo       string  `gorm:"type:text" json:"conteudo"`
	ReferenciaID   *uint   `json:"referencia_id,omitempty"`                  // ID de planta/tarefa relacionada
	ReferenciaTipo *string `gorm:"size:50" json:"referencia_tipo,omitempty"` // "planta", "tarefa", etc.

	// Mídia
	Fotos []Foto `gorm:"polymorphic:Owner;" json:"fotos"`

	// Métricas opcionais
	Clima *ClimaRegistro `gorm:"embedded;embeddedPrefix:clima_" json:"clima,omitempty"`
}

func (r *RegistroDiario) Referencia(db *gorm.DB) interface{} {
	// Retorna a entidade relacionada dinamicamente
	if r.ReferenciaTipo == nil || r.ReferenciaID == nil {
		return nil
	}
	switch *r.ReferenciaTipo {
	case "planta":
		var planta Planta
		db.First(&planta, r.ReferenciaID)
		return planta
	case "tarefa":
		var tarefa Tarefa
		db.First(&tarefa, r.ReferenciaID)
		return tarefa
	}
	return nil
}
