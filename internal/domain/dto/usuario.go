package dto

import "encoding/json"

type UsuarioCreateDTO struct {
    Nome         string `json:"nome" binding:"required,min=3,max=100"`
    Email        string `json:"email" binding:"required,email"`
    Senha        string `json:"senha" binding:"required,min=6"`
    Preferencias json.RawMessage `json:"preferencias"`
}

type UsuarioUpdateDTO struct {
    Nome         string `json:"nome" binding:"min=3,max=100"`
    Email        string `json:"email" binding:"omitempty,email"`
    Preferencias json.RawMessage `json:"preferencias"`
}

type UsuarioResponseDTO struct {
	ID           uint   `json:"id"`
	Nome         string `json:"nome"`
	Email        string `json:"email"`
	Preferencias json.RawMessage `json:"preferencias"`
}
