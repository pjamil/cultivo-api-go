package utils

import "errors"

var (
	ErrInvalidInput     = errors.New("entrada inválida")
	ErrNotFound         = errors.New("recurso não encontrado")
	ErrAlreadyExists    = errors.New("recurso já existe")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrInternalServer   = errors.New("erro interno do servidor")
)
