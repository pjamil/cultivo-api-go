package router

import (
	"github.com/gin-gonic/gin"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/handler"
)

// RegisterDiarioCultivoRoutes registra as rotas de CRUD para a entidade DiarioCultivo.
// Todas as rotas registradas aqui estarão sob a proteção do middleware de autenticação
// aplicado ao grupo de rotas que é passado como argumento.
func RegisterDiarioCultivoRoutes(router *gin.RouterGroup, h *handler.DiarioCultivoHandler) {
	diarios := router.Group("/diarios")
	{
		diarios.POST("", h.Create)
		diarios.GET("", h.GetAllDiarios)
		diarios.GET("/:id", h.GetByID)
		diarios.PUT("/:id", h.Update)
		diarios.DELETE("/:id", h.Delete)
	}
}