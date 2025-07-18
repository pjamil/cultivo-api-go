package entity

import (
	"time"

	"gorm.io/gorm"
)

// RegistroTipo define os tipos de registros que podem ser feitos no diário.
type RegistroTipo string

const (
	RegistroTipoObservacao  RegistroTipo = "observacao"
	RegistroTipoEvento      RegistroTipo = "evento"
	RegistroTipoAprendizado RegistroTipo = "aprendizado"
	RegistroTipoTratamento  RegistroTipo = "tratamento"
	RegistroTipoProblema    RegistroTipo = "problema"
	RegistroTipoColheita    RegistroTipo = "colheita"
	RegistroTipoCrescimento RegistroTipo = "crescimento"
)

// RegistroDiario representa uma entrada no diário de cultivo.
type RegistroDiario struct {
	gorm.Model
	Titulo          string       `gorm:"size:255;not null" json:"titulo"`
	Conteudo        string       `gorm:"type:text;not null" json:"conteudo"`
	Data            time.Time    `gorm:"not null" json:"data"`
	Tipo            RegistroTipo `gorm:"type:varchar(50);not null" json:"tipo"`
	DiarioCultivoID uint         `gorm:"not null" json:"diario_cultivo_id"`
	PlantaID        *uint        `json:"planta_id,omitempty"`

	// Campos polimórficos para associar o registro a uma entidade específica (opcional)
	ReferenciaID   *uint   `json:"referencia_id"`
	ReferenciaTipo *string `gorm:"size:50" json:"referencia_tipo"` // "planta", "tarefa", etc.

	// Mídia
	Fotos []Foto `gorm:"polymorphic:Owner;" json:"fotos"`
}

func (RegistroDiario) TableName() string {
	return "registro_diarios"
}