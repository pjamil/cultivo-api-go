package middleware

import (
	"net/http"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandlerMiddleware captura panics e erros não tratados, retornando uma resposta 500.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Panic recuperado: %v", r)
				utils.RespondWithError(c, http.StatusInternalServerError, "Ocorreu um erro interno no servidor.")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()

		// Se houver erros coletados pelo Gin (ex: de middlewares anteriores ou handlers)
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logrus.WithError(e).Error("Erro não tratado na requisição")
			}
			// Se a resposta já foi escrita, não tente escrever novamente
			if !c.IsAborted() {
				utils.RespondWithError(c, http.StatusInternalServerError, "Ocorreu um erro interno no servidor.")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}
}
