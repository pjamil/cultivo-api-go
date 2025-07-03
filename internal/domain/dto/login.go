package dto

// LoginPayload define a estrutura esperada para o corpo da requisição de login.
type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
