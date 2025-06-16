package models

import (
	"gorm.io/gorm"
)

type ColecaoMidia struct {
	gorm.Model
	DiarioCultivoID uint    `json:"diario_cultivo_id"` // <-- Adicione esta linha
	Nome            string  `gorm:"size:100;not null" json:"nome"`
	Tipo            string  `gorm:"size:50;not null" json:"tipo"` // evolucao, colheita, etc.
	Descricao       string  `gorm:"type:text" json:"descricao"`
	Itens           []Midia `gorm:"foreignKey:ColecaoMidiaID" json:"itens"`
	CapaURL         string  `json:"capa_url"`
}
