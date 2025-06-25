package dto

// DTO para criação de ambiente
type CreateAmbienteDTO struct {
	Nome           string  `json:"nome" binding:"required"`
	Descricao      string  `json:"descricao"`
	Tipo           string  `json:"tipo" binding:"required"`            // Ex: "interno", "externo", "húmido", "seco"
	Comprimento    float64 `json:"comprimento" binding:"required"`     // em centímetros
	Altura         float64 `json:"altura" binding:"required"`          // em centímetros
	Largura        float64 `json:"largura" binding:"required"`         // em centímetros
	TempoExposicao int     `json:"tempo_exposicao" binding:"required"` // em horas
}
