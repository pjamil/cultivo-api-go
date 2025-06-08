package models

import (
	"time"

	"gorm.io/gorm"
)

// TarefaPlanta representa uma tarefa relacionada à planta
type TarefaPlanta struct {
	gorm.Model
	PlantaID      uint       `json:"planta_id"`
	Descricao     string     `gorm:"type:text;not null" json:"descricao"`
	DataPrevista  time.Time  `gorm:"not null" json:"data_prevista"`
	DataRealizada *time.Time `json:"data_realizada,omitempty"`
	Status        string     `gorm:"size:50;not null" json:"status"`
}
