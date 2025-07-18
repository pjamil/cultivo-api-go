package entity

import (
	"time"

	"gorm.io/gorm"
)

type Microclima struct {
	gorm.Model
	AmbienteID   uint      `json:"ambiente_id"`
	DataMedicao  time.Time `json:"data_medicao"`
	Temperatura  float64   `json:"temperatura"`  // °C
	Umidade      float64   `json:"umidade"`      // %
	Luminosidade float64   `json:"luminosidade"` // lux
}
