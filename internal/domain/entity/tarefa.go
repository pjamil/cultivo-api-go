package entity

import (
	"time"

	"gorm.io/gorm"
)



type Tarefa struct {
	gorm.Model
	Tipo           string     `gorm:"size:50;not null" json:"tipo"` // regar, adubar, podar, etc.
	Descricao      string     `gorm:"type:text" json:"descricao"`
	DataAgendada   time.Time  `json:"data_agendada"`
	DataConclusao  *time.Time `json:"data_conclusao"`
	Status         string     `gorm:"size:20;not null" json:"status"` // pendente, concluída, atrasada
	Prioridade     string     `gorm:"size:20" json:"prioridade"`      // baixa, média, alta
	PlantaID       *uint      `json:"planta_id"`                      // NULL para tarefas gerais
	AmbienteID     *uint      `json:"ambiente_id"`
	UsuarioID      uint       `json:"usuario_id"`
	Fotos          []Foto     `gorm:"polymorphic:Owner;" json:"fotos"`
	Recorrente     bool       `json:"recorrente"`
	FrequenciaDias *int       `json:"frequencia_dias"` // NULL para tarefas únicas

	DiarioCultivoID uint `json:"diario_cultivo_id"`
}

type TarefaTemplate struct {
	gorm.Model
	Nome           string `gorm:"size:100;not null" json:"nome"`
	Descricao      string `gorm:"type:text" json:"descricao"`
	Tipo           string `gorm:"size:50;not null" json:"tipo"`
	FrequenciaDias int    `json:"frequencia_dias"`
	Instrucoes     string `gorm:"type:text" json:"instrucoes"`
	EspecieID      *uint  `json:"especie_id"` // NULL para templates gerais
}
