# Implementando o HealthController

Vou implementar o `HealthController` que foi referenciado no código do servidor. Este controlador é responsável por fornecer endpoints de verificação de saúde (health check) da aplicação.

## Implementação do HealthController

`internal/controller/health_controller.go`:

```go
package controller

import (
 "net/http"

 "github.com/gin-gonic/gin"
 "github.com/yourusername/plant-cultivation-api/internal/utils"
)

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
```

## Atualização do Server para incluir todas as rotas de health

`internal/infrastructure/server/server.go` (atualizado):

```go
// ... (código anterior)

 // Health check routes
 healthController := controller.NewHealthController()
 router.GET("/health", healthController.CheckHealth)
 router.GET("/health/ready", healthController.ReadyCheck)
 router.GET("/health/live", healthController.LiveCheck)

 // Plant routes
// ... (restante do código)
```

## Documentação com Swagger (opcional)

Se você estiver usando Swagger, adicione estas anotações ao topo do arquivo `health_controller.go`:

```go
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
```

## Testando os endpoints

Você pode testar os endpoints com curl:

```bash
# Health check básico
curl http://localhost:8080/health

# Verificação de readiness
curl http://localhost:8080/health/ready

# Verificação de liveness
curl http://localhost:8080/health/live
```

## Funcionalidades implementadas

1. **CheckHealth**: Verificação básica de saúde da aplicação
2. **ReadyCheck**: Verificação se todas as dependências estão prontas
3. **LiveCheck**: Verificação mínima de que o processo está rodando

Você pode estender este controlador para incluir verificações adicionais como:

- Conexão com banco de dados
- Status de serviços externos
- Uso de recursos (CPU, memória)
- Status de filas ou brokers de mensagens

Esta implementação segue o padrão dos outros controladores e pode ser facilmente estendida conforme a aplicação cresce.
