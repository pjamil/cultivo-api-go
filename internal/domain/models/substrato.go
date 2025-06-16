package models

import "gorm.io/gorm"

type Substrato struct {
	gorm.Model
	Nome         string  `gorm:"size:100;not null" json:"nome"`
	Composicao   string  `gorm:"type:text" json:"composicao"`
	PH           float64 `json:"ph"`
	RetencaoAgua float64 `json:"retencao_agua"` // 0-1
}
