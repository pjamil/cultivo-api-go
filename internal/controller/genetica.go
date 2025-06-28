package controller

import (
	"errors"
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GeneticaController struct {
	servico service.GeneticaService
}

func NewGeneticaController(servico service.GeneticaService) *GeneticaController {
	return &GeneticaController{servico}
}

// Criar godoc
// @Summary      Cria uma nova genética
// @Description  Cria uma nova genética com os dados fornecidos
// @Tags         genetica
// @Accept       json
// @Produce      json
// @Param        genetica  body      dto.CreateGeneticaDTO  true  "Dados da Genética"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /geneticas [post]
func (ctrl *GeneticaController) Criar(c *gin.Context) {
	var dto dto.CreateGeneticaDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para criar genética")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	geneticaCriada, err := ctrl.servico.Criar(&dto)
	if err != nil {
		logrus.WithError(err).Error("Erro ao criar genética")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar genética"})
		return
	}

	c.JSON(http.StatusCreated, geneticaCriada)
}

// Listar lida com requisições GET para retornar todas as genéticas
func (c *GeneticaController) Listar(ctx *gin.Context) {
	geneticas, err := c.servico.ListarTodas()
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar genéticas")
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, geneticas)
}

// BuscarPorID godoc
// @Summary Busca uma genética por ID
// @Description Retorna informações detalhadas de uma genética
// @Tags genetica
// @Accept json
// @Produce json
// @Param id path int true "Genetica ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /genetica/{id} [get]
func (c *GeneticaController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para buscar genética por ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	genetics, err := c.servico.BuscarPorID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Genética não encontrada ao buscar por ID")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Genética não encontrada"})
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao buscar genética por ID")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar genética"})
		return
	}
	ctx.JSON(http.StatusOK, genetics)
}

// Atualizar godoc
// @Summary      Atualiza uma genética
// @Description  Atualiza uma genética existente com os dados fornecidos
// @Tags         genetica
// @Accept       json
// @Produce      json
// @Param        id        path      int                  true  "ID da Genética"
// @Param        genetica  body      dto.UpdateGeneticaDTO  true  "Dados da Genética para atualização"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /genetica/{id} [put]
func (c *GeneticaController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para atualização de genética")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updateDto dto.UpdateGeneticaDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualização de genética")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	geneticaAtualizada, err := c.servico.Atualizar(uint(id), &updateDto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Genética não encontrada para atualização")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Genética não encontrada"})
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao atualizar genética")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar genética"})
		return
	}

	ctx.JSON(http.StatusOK, geneticaAtualizada)
}
