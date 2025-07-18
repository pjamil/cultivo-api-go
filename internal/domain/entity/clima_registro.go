package entity

type ClimaRegistro struct {
	Temperatura  float64 `json:"temperatura"`  // °C
	Umidade      float64 `json:"umidade"`      // %
	Luminosidade float64 `json:"luminosidade"` // lux
	Precipitacao float64 `json:"precipitacao"` // mm
}
