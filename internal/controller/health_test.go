package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestRouter configura um router Gin com o HealthController para testes.
func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	healthController := NewHealthController(db)
	router.GET("/health", healthController.VerificarSaude)
	router.GET("/health/ready", healthController.VerificarProntidao)
	router.GET("/health/live", healthController.VerificarVitalidade)
	return router
}

// setupTestDB cria uma conexão com um banco de dados SQLite em memória para testes.
func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados de teste: %w", err)
	}
	return db, nil
}

func TestHealthController_VerificarProntidao(t *testing.T) {
	t.Run("Success - Database is connected", func(t *testing.T) {
		// Arrange
		db, err := setupTestDB()
		assert.NoError(t, err)
		router := setupTestRouter(db)

		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/health/ready", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var response StatusResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ready", response.Status)
		assert.Empty(t, response.Error)
	})

	t.Run("Failure - Database is disconnected", func(t *testing.T) {
		// Arrange
		db, err := setupTestDB()
		assert.NoError(t, err)
		router := setupTestRouter(db)

		// Simula a desconexão do banco de dados
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/health/ready", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		var response StatusResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "not ready", response.Status)
		assert.Contains(t, response.Error, "database not available")
	})
}

func TestHealthController_VerificarSaude(t *testing.T) {
	t.Run("Success - Database is connected", func(t *testing.T) {
		// Arrange
		db, err := setupTestDB()
		assert.NoError(t, err)
		router := setupTestRouter(db)

		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response HealthStatusResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response.Status)
		assert.Equal(t, "1.0.0", response.Version)
		assert.Equal(t, "ok", response.Dependencies.Database)
	})

	t.Run("Failure - Database is disconnected", func(t *testing.T) {
		// Arrange
		db, err := setupTestDB()
		assert.NoError(t, err)
		router := setupTestRouter(db)

		// Simula a desconexão do banco de dados
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code) // /health sempre retorna 200

		var response HealthStatusResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response.Status)
		assert.Equal(t, "1.0.0", response.Version)
		assert.Contains(t, response.Dependencies.Database, "error: sql: database is closed")
	})
}

func TestHealthController_VerificarVitalidade(t *testing.T) {
	// Arrange
	router := setupTestRouter(nil) // Não precisa de DB para este teste

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health/live", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response StatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "alive", response.Status)
	assert.Empty(t, response.Error)
}
