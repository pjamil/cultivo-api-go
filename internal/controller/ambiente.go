package controller

import (
	"net/http"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DTO para criação de ambiente
type CreateAmbienteDTO struct {
	Nome      string `json:"nome" binding:"required"`
	Descricao string `json:"descricao"`
}

// Handler para criar novo ambiente
func CreateAmbiente(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreateAmbienteDTO
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ambiente := models.Ambiente{
			Nome:      dto.Nome,
			Descricao: dto.Descricao,
		}

		if err := db.Create(&ambiente).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar ambiente"})
			return
		}

		c.JSON(http.StatusCreated, ambiente)
	}
}
