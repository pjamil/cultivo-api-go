package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DiarioCultivo struct {
	gorm.Model
	Nome       string     `gorm:"size:100;not null" json:"nome"`
	DataInicio time.Time  `json:"data_inicio"`
	DataFim    *time.Time `json:"data_fim,omitempty"`
	Ativo      bool       `gorm:"default:true" json:"ativo"`
	UsuarioID  uint       `json:"usuario_id"`

	// Relacionamentos
	Plantas   []Planta         `gorm:"many2many:diario_plantas;" json:"plantas"`
	Ambientes []Ambiente       `gorm:"many2many:diario_ambientes;" json:"ambientes"`
	Registros []RegistroDiario `json:"registros"`
	Tarefas   []Tarefa         `json:"tarefas"`
	Colecoes  []ColecaoMidia   `json:"colecoes"`

	// Configurações
	Privacidade string `gorm:"size:20;default:'privado'" json:"privacidade"` // privado, publico, compartilhado
	Tags        string `gorm:"size:200" json:"tags"`                         // separadas por vírgula
}

func (d *DiarioCultivo) Resumo() map[string]interface{} {
	return map[string]interface{}{
		"total_plantas": len(d.Plantas),
		"total_tarefas": len(d.Tarefas),
		"periodo": fmt.Sprintf("%s a %s",
			d.DataInicio.Format("02/01/2006"),
			ternary(d.DataFim != nil,
				d.DataFim.Format("02/01/2006"),
				"presente")),
		"ultimo_registro": d.ultimoRegistro(),
	}
}

// func (d *DiarioCultivo) Timeline() []TimelineItem {
// 	// Combina registros, tarefas e eventos em ordem cronológica
// 	// Implementação personalizável
// }

// ultimoRegistro retorna o último registro do diário, se existir.
func (d *DiarioCultivo) ultimoRegistro() interface{} {
	if len(d.Registros) == 0 {
		return nil
	}
	// Supondo que Registros está ordenado por data crescente
	return d.Registros[len(d.Registros)-1]
}

// ternary is a helper function to mimic the ternary operator for string values.
func ternary(condition bool, a, b string) string {
	if condition {
		return a
	}
	return b
}
