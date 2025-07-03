package middleware

import (
	"net/http"
	"strings"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware verifica a validade do token JWT em cada requisição.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Warn("Requisição sem cabeçalho de autorização")
			utils.RespondWithError(c, http.StatusUnauthorized, "Token de autenticação ausente")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logrus.Warnf("Formato de cabeçalho de autorização inválido: %s", authHeader)
			utils.RespondWithError(c, http.StatusUnauthorized, "Formato de token inválido. Use Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			logrus.WithError(err).Warnf("Token inválido ou expirado: %s", tokenString)
			utils.RespondWithError(c, http.StatusUnauthorized, "Token inválido ou expirado")
			c.Abort()
			return
		}

		// Adiciona o ID do usuário ao contexto da requisição para uso posterior nos handlers
		c.Set("userID", userID)
		c.Next()
	}
}
