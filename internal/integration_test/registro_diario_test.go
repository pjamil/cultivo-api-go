package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/stretchr/testify/assert" // Importe o pacote utils
	"gorm.io/gorm"
)

// createTestUserAndDiario é uma função auxiliar para criar um usuário e um diário de cultivo para os testes.
func createTestUserAndDiario(t *testing.T, db *gorm.DB) (user entity.Usuario, diario entity.DiarioCultivo, token string) {
	// Criar usuário
	user = entity.Usuario{
		Nome:      "Usuário de Teste Diário",
		Email:     fmt.Sprintf("diario.user.%d@example.com", time.Now().UnixNano()),
		SenhaHash: "testpassword", // Use a senha real para o hash
	}
	hashedPassword, err := utils.HashPassword(user.SenhaHash)
	assert.NoError(t, err)
	user.SenhaHash = hashedPassword
	err = db.Create(&user).Error
	assert.NoError(t, err)

	// Criar diário
	diario = entity.DiarioCultivo{
		Nome:      "Diário de Teste para Registros",
		UsuarioID: user.ID,
	}
	err = db.Create(&diario).Error
	assert.NoError(t, err)

	// Gerar token
	token, err = utils.GenerateToken(user.ID)
	assert.NoError(t, err)

	return user, diario, token
}

func TestCriarRegistroDiario(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB
	LimparBancoDeDados(db)

	_, diario, token := createTestUserAndDiario(t, db)

	t.Run("Success - 201 Created", func(t *testing.T) {
		// Arrange
		payload := dto.CreateRegistroDiarioDTO{
			Titulo:   "Primeira Observação",
			Conteudo: "A planta está crescendo bem.",
			Data:     time.Now(),
			Tipo:     entity.RegistroTipoObservacao,
		}
		body, _ := json.Marshal(payload)

		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/diarios-cultivo/%d/registros", diario.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		fmt.Println("Response Body:", w.Body.String())
		var response dto.RegistroDiarioResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, payload.Titulo, response.Titulo)
		assert.NotZero(t, response.ID)
	})

	t.Run("Failure - 400 Bad Request (Invalid Payload)", func(t *testing.T) {
		// Arrange
		payload := dto.CreateRegistroDiarioDTO{
			Titulo: "Título Curto", // Conteúdo é obrigatório
		}
		body, _ := json.Marshal(payload)

		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/diarios-cultivo/%d/registros", diario.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failure - 404 Not Found (Diario does not exist)", func(t *testing.T) {
		// Arrange
		payload := dto.CreateRegistroDiarioDTO{
			Titulo:   "Observação em Diário Inexistente",
			Conteudo: "Este teste deve falhar.",
			Data:     time.Now(),
			Tipo:     entity.RegistroTipoObservacao,
		}
		body, _ := json.Marshal(payload)

		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/diarios-cultivo/99999/registros", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Failure - 401 Unauthorized", func(t *testing.T) {
		// Arrange
		payload := dto.CreateRegistroDiarioDTO{
			Titulo:   "Teste sem token",
			Conteudo: "Deve falhar.",
			Data:     time.Now(),
			Tipo:     entity.RegistroTipoObservacao,
		}
		body, _ := json.Marshal(payload)

		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/diarios-cultivo/%d/registros", diario.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestListarRegistrosDiario(t *testing.T) {
	db := GetTestDB().DB
	router := GetTestRouter()

	_, diario, token := createTestUserAndDiario(t, db)

	// Criar alguns registros para o teste
	for i := 0; i < 3; i++ {
		registro := entity.RegistroDiario{
			Titulo:          fmt.Sprintf("Registro de Teste %d", i+1),
			Conteudo:        "Conteúdo do registro.",
			Data:            time.Now(),
			Tipo:            entity.RegistroTipoEvento,
			DiarioCultivoID: diario.ID,
		}
		err := db.Create(&registro).Error
		assert.NoError(t, err)
	}

	t.Run("Success - 200 OK", func(t *testing.T) {
		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/diarios-cultivo/%d/registros", diario.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), response.Total)

		var registros []dto.RegistroDiarioResponseDTO
		if response.Data == nil {
			response.Data = []byte("[]")
		}
		err = json.Unmarshal(response.Data, &registros)
		assert.NoError(t, err)
		assert.Len(t, registros, 3)
		assert.Equal(t, "Registro de Teste 3", registros[0].Titulo) // Ordenado por data desc
	})
}
