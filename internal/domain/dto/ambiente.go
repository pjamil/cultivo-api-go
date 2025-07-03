package dto

// DTO para criação de ambiente
type CreateAmbienteDTO struct {
	Nome           string  `json:"nome" binding:"required"`
	Descricao      string  `json:"descricao"`
	Tipo           string  `json:"tipo" binding:"required,oneof='interno' 'externo' 'húmido' 'seco'"`            // Ex: "interno", "externo", "húmido", "seco"
	Comprimento    float64 `json:"comprimento" binding:"required,gt=0"`     // em centímetros
	Altura         float64 `json:"altura" binding:"required,gt=0"`          // em centímetros
	Largura        float64 `json:"largura" binding:"required,gt=0"`         // em centímetros
	TempoExposicao int     `json:"tempo_exposicao" binding:"required,gt=0"` // em horas
}

// DTO para atualização de ambiente
type UpdateAmbienteDTO struct {
	Nome           string  `json:"nome"`
	Descricao      string  `json:"descricao"`
	Tipo           string  `json:"tipo"`
	Comprimento    float64 `json:"comprimento"`
	Altura         float64 `json:"altura"`
	Largura        float64 `json:"largura"`
	TempoExposicao int     `json:"tempo_exposicao"`
}

// DTO para resposta de ambiente
type AmbienteResponseDTO struct {
	ID             uint    `json:"id"`
	Nome           string  `json:"nome"`
	Descricao      string  `json:"descricao,omitempty"`
	Tipo           string  `json:"tipo"`
	Comprimento    float64 `json:"comprimento"`
	Altura         float64 `json:"altura"`
	Largura        float64 `json:"largura"`
	TempoExposicao int     `json:"tempo_exposicao"`
	Orientacao     string  `json:"orientacao"`
}
