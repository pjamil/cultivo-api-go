package utils

import "errors"

// Erros comuns da aplicação
var (
	ErrNotFound         = errors.New("recurso não encontrado")
	ErrInvalidInput     = errors.New("entrada inválida")
	ErrUnauthorized     = errors.New("não autorizado")
	ErrConflict         = errors.New("conflito de dados")
	ErrInternalServer   = errors.New("erro interno do servidor")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
)

// NewError cria um novo erro com uma mensagem formatada.
func NewError(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
