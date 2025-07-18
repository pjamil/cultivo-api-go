package middleware

import (
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware retorna um middleware Gin que adiciona um logger contextualizado
// com um request_id único para cada requisição.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Inicia o timer
		start := time.Now()

		// Gera um request_id único para esta requisição
		// TODO: Implementar um gerador de UUID mais robusto se necessário
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = "req-" + strconv.FormatInt(time.Now().UnixNano(), 10)
		}

		// Cria um logger de entrada com campos padrão para esta requisição
		logger := logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
		})

		// Adiciona o logger ao contexto do Gin para que possa ser acessado nas camadas inferiores
		c.Set("logger", logger)

		// Processa a requisição
		c.Next()

		// Calcula a duração da requisição
		duration := time.Since(start)

		// Adiciona informações de resposta ao log
		logger.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"duration": duration.String(),
			"size":     c.Writer.Size(),
		}).Info("Requisição processada")
	}
}
