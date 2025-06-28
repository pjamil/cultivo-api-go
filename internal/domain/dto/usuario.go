package dto

type UsuarioCreateDTO struct {
    Nome         string `json:"nome" binding:"required"`
    Email        string `json:"email" binding:"required,email"`
    Senha        string `json:"senha" binding:"required,min=6"`
    Preferencias string `json:"preferencias"`
}

type UsuarioUpdateDTO struct {
    Nome         string `json:"nome"`
    Preferencias string `json:"preferencias"`
}

type UsuarioResponseDTO struct {
	ID           uint   `json:"id"`
	Nome         string `json:"nome"`
	Email        string `json:"email"`
	Preferencias string `json:"preferencias"`
}
