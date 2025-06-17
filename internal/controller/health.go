package controller

import (
	"net/http"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
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
	// Você pode adicionar dependências aqui como banco de dados, serviços, etc.
	// Por exemplo: db *gorm.DB para verificar a conexão com o banco
}

func NewHealthController() *HealthController {
	return &HealthController{
		// Inicialize quaisquer dependências aqui
	}
}

// CheckHealth godoc
// @Summary Verifica a saúde da aplicação
// @Description Retorna o status da API e seus componentes
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "status": "ok"
// @Router /health [get]
func (c *HealthController) CheckHealth(ctx *gin.Context) {
	// Você pode adicionar verificações adicionais aqui, como:
	// - Conexão com banco de dados
	// - Conexão com serviços externos
	// - Uso de recursos

	healthStatus := gin.H{
		"status":  "ok",
		"version": "1.0.0",
		// Adicione mais informações de saúde conforme necessário
	}

	utils.RespondWithJSON(ctx, http.StatusOK, healthStatus)
}

// ReadyCheck godoc
// @Summary Verifica se a aplicação está pronta para receber tráfego
// @Description Verifica todas as dependências necessárias para a aplicação funcionar
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "status": "ready"
// @Failure 503 {object} map[string]interface{} "status": "not ready"
// @Router /health/ready [get]
func (c *HealthController) ReadyCheck(ctx *gin.Context) {
	// Implemente verificações mais rigorosas aqui
	// Por exemplo:
	// if err := c.db.Exec("SELECT 1").Error; err != nil {
	//     utils.RespondWithJSON(ctx, http.StatusServiceUnavailable, gin.H{
	//         "status": "not ready",
	//         "error":  "database not available",
	//     })
	//     return
	// }

	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{
		"status":  "ready",
		"version": "1.0.0",
	})
}

// LiveCheck godoc
// @Summary Verifica se a aplicação está viva
// @Description Verificação simples de que o processo está em execução
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "status": "alive"
// @Router /health/live [get]
func (c *HealthController) LiveCheck(ctx *gin.Context) {
	// Verificação mínima de que o processo está rodando
	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{
		"status": "alive",
	})
}

func (c *HealthController) Check(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status":  "up",
		"version": "1.0.0",
	})
}
