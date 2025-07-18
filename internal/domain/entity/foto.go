package entity

import (
	"gorm.io/gorm"
)

// Foto representa uma foto no sistema
type Foto struct {
	gorm.Model
	OwnerID    uint   `gorm:"index" json:"owner_id"`
	OwnerType  string `gorm:"size:50;index" json:"owner_type"` // "planta", "tarefa", etc.
	URL        string `gorm:"not null" json:"url"`
	Descricao  string `gorm:"type:text" json:"descricao"`
	UsuarioID  uint   `json:"usuario_id"`
	AmbienteID *uint  `json:"ambiente_id,omitempty"` // <-- Adicione esta linha
}


