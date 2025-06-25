package dto

type CreateMeioCultivoDTO struct {
	Tipo      string `json:"tipo" binding:"required"`
	Descricao string `json:"descricao"`
}
