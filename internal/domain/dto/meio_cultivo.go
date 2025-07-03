package dto

type CreateMeioCultivoDTO struct {
	Tipo      string `json:"tipo" binding:"required,oneof='solo' 'hidroponia' 'coco' 'lã de rocha' 'turfa'"`
	Descricao string `json:"descricao" binding:"max=255"`
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
