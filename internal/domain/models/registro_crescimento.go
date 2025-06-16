package models

import (
	"time"

	"gorm.io/gorm"
)

type Planta2 struct {
	gorm.Model
	NomePopular          string                `gorm:"size:100;not null" json:"nome_popular"`
	NomeCientifico       string                `gorm:"size:100" json:"nome_cientifico"`
	EspecieID            uint                  `json:"especie_id"`
	DataPlantio          time.Time             `json:"data_plantio"`
	DataGerminacao       *time.Time            `json:"data_germinacao"`
	DataColheita         *time.Time            `json:"data_colheita"`
	Status               string                `gorm:"size:50;not null" json:"status"` // ativa, inativa, colhida, etc.
	AmbienteID           uint                  `json:"ambiente_id"`
	VasoID               *uint                 `json:"vaso_id"`
	ProprietarioID       uint                  `json:"proprietario_id"`
	RegistrosCrescimento []RegistroCrescimento `json:"registros_crescimento"`
	Tarefas              []Tarefa              `json:"tarefas"`
	Fotos                []Foto                `gorm:"polymorphic:Owner;" json:"fotos"`
	Observacoes          string                `gorm:"type:text" json:"observacoes"`
}

type RegistroCrescimento struct {
	gorm.Model
	PlantaID       uint    `json:"planta_id"`
	Altura         float64 `json:"altura"` // em cm
	NumeroFolhas   int     `json:"numero_folhas"`
	DiamentroCaule float64 `json:"diametro_caule"`         // em mm
	Estagio        string  `gorm:"size:50" json:"estagio"` // germinação, vegetativo, etc.
	Observacoes    string  `gorm:"type:text" json:"observacoes"`
}
