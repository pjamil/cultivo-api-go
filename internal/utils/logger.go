package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetLoggerFromContext recupera a instância do logger com campos (como request_id) do contexto do Gin.
// Se nenhum logger for encontrado por algum motivo, ele retorna um logger padrão para garantir
// que a aplicação não quebre, embora o contexto de rastreamento seja perdido.
func GetLoggerFromContext(c *gin.Context) *logrus.Entry {
	logger, exists := c.Get("logger")
	if !exists {
		return logrus.WithField("context", "missing")
	}

	if l, ok := logger.(*logrus.Entry); ok {
		return l
	}

	return logrus.WithField("context", "invalid_type")
}
