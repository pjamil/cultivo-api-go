package dto

import "encoding/json"

// PaginationParams define os parâmetros de paginação para requisições.
type PaginationParams struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

// PaginatedResponse define a estrutura de resposta para dados paginados.
type PaginatedResponse struct {
	Data  json.RawMessage `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}
