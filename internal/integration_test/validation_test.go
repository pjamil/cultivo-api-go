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
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	tu "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateWithInvalidData(t *testing.T) {
	router := GetTestRouter()

	// Cenário 1: Criar Usuário com Email Inválido
	invalidUserPayload := dto.UsuarioCreateDTO{
		Nome:         "Invalid User",
		Email:        "invalid-email", // Email inválido
		Senha:        "password123",
		Preferencias: "",
	}
	jsonPayload, _ := json.Marshal(invalidUserPayload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/usuarios", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Email") // Verifica se a mensagem de erro contém o campo Email

	// Cenário 2: Criar Planta com Nome Vazio (campo obrigatório)
	plantioTime, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	assert.NoError(t, err)
	invalidPlantaPayload := dto.CreatePlantaDTO{
		Nome:        "", // Nome vazio
		Especie:     "Cannabis Sativa",
								DataPlantio: plantioTime,
	}
	jsonPayload, _ = json.Marshal(invalidPlantaPayload)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/plantas", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Nome") // Verifica se a mensagem de erro contém o campo Nome

	// Cenário 3: Criar Ambiente com Tipo Inválido (oneof)
	invalidAmbientePayload := dto.CreateAmbienteDTO{
		Nome:           "Ambiente Teste",
		Tipo:           "tipo-invalido", // Tipo inválido
		Comprimento:    100,
		Altura:         50,
		Largura:        50,
		TempoExposicao: 12,
	}
	jsonPayload, _ = json.Marshal(invalidAmbientePayload)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/ambientes", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Tipo") // Verifica se a mensagem de erro contém o campo Tipo
}

func TestUpdateWithInvalidData(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB

	// Limpar o banco de dados antes do teste
	db.Exec("DELETE FROM usuarios")

	// 1. Criar um usuário válido para atualização
	password := "validpassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	validUser := models.Usuario{
		Nome:         "Valid User",
		Email:        "valid@example.com",
		SenhaHash:    hashedPassword,
		Preferencias: "",
	}
	err = db.Create(&validUser).Error
	assert.NoError(t, err)

	// Cenário 1: Atualizar Usuário com Email Inválido
	invalidUpdatePayload := dto.UsuarioUpdateDTO{
		Nome:  "Updated User",
		Email: "invalid-email", // Email inválido
	}
	jsonPayload, _ := json.Marshal(invalidUpdatePayload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/usuarios/%d", validUser.ID), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Email")

	// Cenário 2: Atualizar Planta com Nome muito curto
	// Primeiro, crie uma planta válida
	planta := models.Planta{
		Nome:        "Planta Original",
		Especie:     "Especie Teste",
				DataPlantio: tu.TimePtr(time.Now()),
		Status:      "vegetativo",
		UsuarioID:   validUser.ID,
	}
	err = db.Create(&planta).Error
	assert.NoError(t, err)

	invalidPlantaUpdatePayload := dto.UpdatePlantaDTO{
		Nome: "ab", // Nome muito curto (min=3)
	}
	jsonPayload, _ = json.Marshal(invalidPlantaUpdatePayload)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/plantas/%d", planta.ID), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Nome")
}