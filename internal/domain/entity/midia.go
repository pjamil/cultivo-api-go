package entity

import (
	"time"

	"gorm.io/gorm"
)

type Midia struct {
	gorm.Model
	Tipo           string    `gorm:"size:20;not null" json:"tipo"` // foto, video, audio
	URL            string    `gorm:"not null" json:"url"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	DataCaptura    time.Time `json:"data_captura"`
	AutorID        uint      `json:"autor_id"`
	Descricao      string    `gorm:"type:text" json:"descricao"`
	Coordenadas    string    `gorm:"size:50" json:"coordenadas"` // formato "lat,long"
	ColecaoMidiaID uint      `json:"colecao_midia_id"`           // ID da coleção de mídia associada
}
