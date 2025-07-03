package dto

type UsuarioCreateDTO struct {
    Nome         string `json:"nome" binding:"required,min=3,max=100"`
    Email        string `json:"email" binding:"required,email"`
    Senha        string `json:"senha" binding:"required,min=6"`
    Preferencias string `json:"preferencias" binding:"max=255"`
}

type UsuarioUpdateDTO struct {
    Nome         string `json:"nome" binding:"min=3,max=100"`
    Preferencias string `json:"preferencias" binding:"max=255"`
}

type UsuarioResponseDTO struct {
	ID           uint   `json:"id"`
	Nome         string `json:"nome"`
	Email        string `json:"email"`
	Preferencias string `json:"preferencias"`
}
