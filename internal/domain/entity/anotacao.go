package entity

import "gorm.io/gorm"

type Anotacao struct {
	gorm.Model
	Titulo     string `gorm:"size:100;not null" json:"titulo"`
	Conteudo   string `gorm:"type:text;not null" json:"conteudo"`
	UsuarioID  uint   `json:"usuario_id"`
	PlantaID   *uint  `json:"planta_id"`
	AmbienteID *uint  `json:"ambiente_id"`
	TarefaID   *uint  `json:"tarefa_id"`
	Tags       string `gorm:"size:200" json:"tags"` // separadas por vírgula
}
