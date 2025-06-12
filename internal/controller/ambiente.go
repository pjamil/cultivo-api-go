package controller

import (
	"net/http"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DTO para criação de ambiente
type CreateAmbienteDTO struct {
	Nome           string  `json:"nome" binding:"required"`
	Descricao      string  `json:"descricao"`
	Tipo           string  `json:"tipo" binding:"required"`            // Ex: "interno", "externo", "húmido", "seco"
	Comprimento    float64 `json:"comprimento" binding:"required"`     // em centímetros
	Altura         float64 `json:"altura" binding:"required"`          // em centímetros
	Largura        float64 `json:"largura" binding:"required"`         // em centímetros
	TempoExposicao int     `json:"tempo_exposicao" binding:"required"` // em horas
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
			Nome:           dto.Nome,
			Descricao:      dto.Descricao,
			Tipo:           dto.Tipo,
			Comprimento:    dto.Comprimento,
			Altura:         dto.Altura,
			Largura:        dto.Largura,
			TempoExposicao: dto.TempoExposicao,
		}

		if err := db.Create(&ambiente).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar ambiente"})
			return
		}

		c.JSON(http.StatusCreated, ambiente)
	}
}

// Handler para listar todos os ambientes
func ListAmbientes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ambientes []models.Ambiente
		if err := db.Find(&ambientes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar ambientes"})
			return
		}

		c.JSON(http.StatusOK, ambientes)
	}
}
