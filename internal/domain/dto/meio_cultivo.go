package dto

type CreateMeioCultivoDTO struct {
	Tipo      string `json:"tipo" binding:"required"`
	Descricao string `json:"descricao"`
}

type UpdateMeioCultivoDTO struct {
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type MeioCultivoResponseDTO struct {
	ID        uint   `json:"id"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}
