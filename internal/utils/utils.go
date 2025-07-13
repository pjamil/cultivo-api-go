package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents a standardized error response format.
type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"` // Can be a string, map[string]string, or any other relevant error details
}

func RespondWithError(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, ErrorResponse{Message: message, Details: details})
}

func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

// Common API errors
var (
	ErrInvalidInput     = errors.New("entrada inválida")
	ErrNotFound         = errors.New("recurso não encontrado")
	ErrInternalServer   = errors.New("erro interno do servidor")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
)

// GetErrorMsg returns a user-friendly error message for validation errors.
func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Este campo é obrigatório"
	case "lte":
		return "Deve ser menor ou igual a " + fe.Param()
	case "gte":
		return "Deve ser maior ou igual a " + fe.Param()
	case "email":
		return "Formato de e-mail inválido"
	case "min":
		return "Deve ter no mínimo " + fe.Param() + " caracteres"
	case "max":
		return "Deve ter no máximo " + fe.Param() + " caracteres"
	case "alphanum":
		return "Deve conter apenas caracteres alfanuméricos"
	case "numeric":
		return "Deve conter apenas números"
	case "len":
		return "Deve ter exatamente " + fe.Param() + " caracteres"
	case "oneof":
		return "Deve ser um dos valores: " + fe.Param()
	}
	return "Erro de validação desconhecido"
}