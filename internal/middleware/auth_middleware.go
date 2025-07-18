package middleware

import (
	"net/http"
	"strings"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware verifica a validade do token JWT em cada requisição.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		requestID := uuid.New().String()

		log := logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"path":       c.Request.URL.Path,
		})

		log.Infof("AuthMiddleware: Authorization Header: %s", authHeader)

		if authHeader == "" {
			log.Warn("Requisição sem cabeçalho de autorização")
			utils.RespondWithError(c, http.StatusUnauthorized, "Token de autenticação ausente", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn("Formato de cabeçalho de autorização inválido")
			utils.RespondWithError(c, http.StatusUnauthorized, "Formato de token inválido. Use Bearer <token>", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		log.Infof("AuthMiddleware: Token String: %s", tokenString)
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			log.WithError(err).Warn("Falha na validação do token")
			utils.RespondWithError(c, http.StatusUnauthorized, "Token inválido ou expirado", nil)
			c.Abort()
			return
		}

		// Adiciona o ID do usuário ao contexto da requisição para uso posterior nos handlers
		log.WithField("user_id", userID).Info("Usuário autenticado com sucesso")
		c.Set("userID", userID)

		// Adiciona o logger com campos ao contexto para ser usado nos controllers
		c.Set("logger", log)

		c.Next()
	}
}
