package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
)

func TestLoginSuccess(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB

	// Limpar o banco de dados antes do teste
	db.Exec("DELETE FROM usuarios")

	// 1. Criar um usuário de teste
	password := "testpassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	testUser := models.Usuario{
		Nome:        "Test User",
		Email:       "test@example.com",
		SenhaHash:   hashedPassword,
		Preferencias: "{}",
	}
	err = db.Create(&testUser).Error
	assert.NoError(t, err)

	// 2. Preparar o payload de login
	loginPayload := dto.LoginPayload{
		Email:    "test@example.com",
		Password: password,
	}
	jsonPayload, _ := json.Marshal(loginPayload)

	// 3. Enviar a requisição POST para /api/v1/login
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 4. Verificar a resposta
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
	assert.NotEmpty(t, response["token"])

	// Opcional: Validar o token (apenas para garantir que é um JWT válido)
	userID, err := utils.ValidateToken(response["token"])
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprint(testUser.ID), userID)
}

func TestLoginInvalidCredentials(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB

	// Limpar o banco de dados antes do teste
	db.Exec("DELETE FROM usuarios")

	// Criar um usuário de teste para o cenário de senha incorreta
	password := "correctpassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	testUser := models.Usuario{
		Nome:        "Another User",
		Email:       "another@example.com",
		SenhaHash:   hashedPassword,
		Preferencias: "{}",
	}
	err = db.Create(&testUser).Error
	assert.NoError(t, err)

	// Cenário 1: Email não registrado
	loginPayloadInvalidEmail := dto.LoginPayload{
		Email:    "nonexistent@example.com",
		Password: "anypassword",
	}
	jsonPayloadInvalidEmail, _ := json.Marshal(loginPayloadInvalidEmail)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonPayloadInvalidEmail))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusUnauthorized, w1.Code)
	assert.Contains(t, w1.Body.String(), "Credenciais inválidas")

	// Cenário 2: Senha incorreta
	loginPayloadWrongPassword := dto.LoginPayload{
		Email:    "another@example.com",
		Password: "wrongpassword",
	}
	jsonPayloadWrongPassword, _ := json.Marshal(loginPayloadWrongPassword)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonPayloadWrongPassword))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusUnauthorized, w2.Code)
	assert.Contains(t, w2.Body.String(), "Credenciais inválidas")
}

func TestProtectedRouteWithValidToken(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB

	// Limpar o banco de dados antes do teste
	db.Exec("DELETE FROM usuarios")

	// 1. Criar um usuário de teste
	password := "securepassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	testUser := models.Usuario{
		Nome:        "Protected User",
		Email:       "protected@example.com",
		SenhaHash:   hashedPassword,
		Preferencias: "{}",
	}
	err = db.Create(&testUser).Error
	assert.NoError(t, err)

	// 2. Realizar o login para obter um token JWT válido
	loginPayload := dto.LoginPayload{
		Email:    "protected@example.com",
		Password: password,
	}
	jsonPayload, _ := json.Marshal(loginPayload)

	wLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonPayload))
	reqLogin.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wLogin, reqLogin)

	assert.Equal(t, http.StatusOK, wLogin.Code)

	var loginResponse map[string]string
	err = json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	assert.NoError(t, err)
	token := loginResponse["token"]

	// 3. Fazer uma requisição para uma rota protegida (ex: GET /api/v1/plantas)
	wProtected := httptest.NewRecorder()
	reqProtected, _ := http.NewRequest("GET", "/api/v1/plantas", nil)
	reqProtected.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(wProtected, reqProtected)

	// 4. Verificar se a resposta é 200 OK
	assert.Equal(t, http.StatusOK, wProtected.Code)
}

func TestProtectedRouteWithoutToken(t *testing.T) {
	router := GetTestRouter()

	// 1. Fazer uma requisição para uma rota protegida (ex: GET /api/v1/plantas) sem token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/plantas", nil)
	router.ServeHTTP(w, req)

	// 2. Verificar se a resposta é 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token de autenticação ausente")
}

func TestProtectedRouteWithInvalidToken(t *testing.T) {
	router := GetTestRouter()

	// 1. Fazer uma requisição para uma rota protegida (ex: GET /api/v1/plantas) com token inválido
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/plantas", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.string")
	router.ServeHTTP(w, req)

	// 2. Verificar se a resposta é 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token inválido ou expirado")
}
