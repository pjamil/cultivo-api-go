package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Erros comuns da aplicação
var (
	ErrNotFound           = errors.New("recurso não encontrado")
	ErrInvalidInput       = errors.New("entrada inválida")
	ErrUnauthorized       = errors.New("não autorizado")
	ErrConflict           = errors.New("conflito de dados")
	ErrInternalServer     = errors.New("erro interno do servidor")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
)

// NewError cria um novo erro com uma mensagem formatada.
func NewError(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// GetErrorMsg retorna uma mensagem de erro amigável para um erro de validação.
func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("O campo %s é obrigatório", fe.Field())
	case "email":
		return fmt.Sprintf("O campo %s deve ser um email válido", fe.Field())
	case "min":
		return fmt.Sprintf("O campo %s deve ter no mínimo %s caracteres", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("O campo %s deve ter no máximo %s caracteres", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("O campo %s deve ser um dos seguintes valores: %s", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("O campo %s deve ser maior que %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("O campo %s deve ser menor ou igual a %s", fe.Field(), fe.Param())
	}
	return fe.Error()
}