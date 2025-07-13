package controller

import (
	"context"
	"net/http"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @title Plant Cultivation API
// @version 1.0
// @description API para gerenciamento de cultivo de plantas
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
type HealthController struct {
	db *gorm.DB
}

// HealthStatusResponse representa a resposta do endpoint de saúde.
type HealthStatusResponse struct {
	Status       string             `json:"status" example:"ok"`
	Version      string             `json:"version" example:"1.0.0"`
	Dependencies HealthDependencies `json:"dependencies"`
}

// HealthDependencies representa o status das dependências.
type HealthDependencies struct {
	Database string `json:"database" example:"ok"`
}

// StatusResponse representa uma resposta de status simples para os endpoints de prontidão e vitalidade.
type StatusResponse struct {
	Status string `json:"status" example:"ready"`
	Error  string `json:"error,omitempty" example:"database not available"`
}

func NewHealthController(db *gorm.DB) *HealthController {
	return &HealthController{
		db: db,
	}
}

// CheckHealth godoc
// @Summary Verifica a saúde da aplicação
// @Description Retorna o status da API e seus componentes
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} controller.HealthStatusResponse "API e dependências estão saudáveis"
// @Router /health [get]
func (c *HealthController) VerificarSaude(ctx *gin.Context) {
	dbStatus := "ok"
	sqlDB, err := c.db.DB()
	if err != nil {
		dbStatus = "error: failed to get db instance"
	} else {
		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 1*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(timeoutCtx); err != nil {
			dbStatus = "error: " + err.Error()
		}
	}

	healthStatus := HealthStatusResponse{
		Status:  "ok",
		Version: "1.0.0",
		Dependencies: HealthDependencies{
			Database: dbStatus,
		},
	}

	utils.RespondWithJSON(ctx, http.StatusOK, &healthStatus)
}

// ReadyCheck godoc
// @Summary Verifica se a aplicação está pronta para receber tráfego
// @Description Verifica todas as dependências necessárias para a aplicação funcionar
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} controller.StatusResponse "Aplicação está pronta"
// @Failure 503 {object} controller.StatusResponse "Aplicação não está pronta"
// @Router /health/ready [get]
func (c *HealthController) VerificarProntidao(ctx *gin.Context) {
	sqlDB, err := c.db.DB()
	if err != nil {
		utils.RespondWithJSON(ctx, http.StatusServiceUnavailable, StatusResponse{Status: "not ready", Error: "failed to get db instance"})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 1*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(timeoutCtx); err != nil {
		utils.RespondWithJSON(ctx, http.StatusServiceUnavailable, StatusResponse{Status: "not ready", Error: "database not available: " + err.Error()})
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, StatusResponse{
		Status: "ready",
	})
}

// LiveCheck godoc
// @Summary Verifica se a aplicação está viva
// @Description Verificação simples de que o processo está em execução
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} controller.StatusResponse "Aplicação está viva"
// @Router /health/live [get]
func (c *HealthController) VerificarVitalidade(ctx *gin.Context) {
	// Verificação mínima de que o processo está rodando
	utils.RespondWithJSON(ctx, http.StatusOK, StatusResponse{
		Status: "alive",
	})
}
