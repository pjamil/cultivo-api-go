package entity

import (
	"time"

	"gorm.io/gorm"
)

// DiarioCultivo representa o modelo de domínio para um diário de cultivo.
type DiarioCultivo struct {
	gorm.Model
	Nome        string    `gorm:"type:varchar(255);not null"`
	DataInicio  time.Time `gorm:"not null"`
	DataFim     *time.Time
	Privacidade string    `gorm:"type:varchar(50);not null"` // Ex: "publico", "privado"
	Tags        string    `gorm:"type:varchar(255)"`
	UsuarioID   uint      `gorm:"not null"`
	// Relacionamentos
	Usuario   Usuario
	Plantas   []Planta   `gorm:"many2many:diario_cultivo_plantas;"`
	Ambientes []Ambiente `gorm:"many2many:diario_cultivo_ambientes;"`
	Registros []RegistroDiario
}