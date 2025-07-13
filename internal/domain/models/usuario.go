package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Usuario representa um usuário do sistema de jardinagem
type UsuarioCreateDTO struct {
	Nome         string `json:"nome" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Senha        string `json:"senha" binding:"required,min=6"`
	Preferencias json.RawMessage `json:"preferencias"`
}

type Usuario struct {
	gorm.Model
	Nome         string   `gorm:"size:100;not null" json:"nome"`
	Email        string   `gorm:"size:100;unique;not null" json:"email"`
	SenhaHash    string   `gorm:"size:255;not null" json:"-"`
	Preferencias json.RawMessage `gorm:"type:json" json:"preferencias"` // Configurações em JSON
	Plantas      []Planta `gorm:"foreignKey:UsuarioID" json:"plantas"`
	Tarefas      []Tarefa `gorm:"foreignKey:UsuarioID" json:"tarefas"`
}

type Lembrete struct {
	gorm.Model
	UsuarioID  uint      `json:"usuario_id"`
	Mensagem   string    `gorm:"type:text;not null" json:"mensagem"`
	DataHora   time.Time `json:"data_hora"`
	Repetir    bool      `json:"repetir"`
	Frequencia string    `gorm:"size:20" json:"frequencia"` // diária, semanal, etc.
	Lido       bool      `json:"lido"`
}
