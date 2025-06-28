package controller

import (
	"errors"
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AmbienteController struct {
	servico service.AmbienteService
}

func NewAmbienteController(servico service.AmbienteService) *AmbienteController {
	return &AmbienteController{servico}
}

// Criar godoc
// @Summary      Cria um novo ambiente
// @Description  Cria um novo ambiente com os dados fornecidos
// @Tags         ambiente
// @Accept       json
// @Produce      json
// @Param        ambiente  body      dto.CreateAmbienteDTO  true  "Dados do Ambiente"
// @Success      201      {object}  models.Ambiente
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /ambiente [post]
func (c *AmbienteController) Criar(ctx *gin.Context) {
	var dto dto.CreateAmbienteDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para criar ambiente")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ambiente, err := c.servico.Criar(&dto)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "create_ambiente",
			"error":     err,
		}).Error("Erro ao criar ambiente")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar ambiente"})
		return
	}
	ctx.JSON(http.StatusCreated, ambiente)
}

// Listar godoc
// @Summary      Lista todos os ambientes
// @Description  Retorna uma lista de todos os ambientes cadastrados
// @Tags         ambiente
// @Produce      json
// @Success      200  {array}   models.Ambiente
// @Failure      500  {object}  map[string]string
// @Router       /ambiente [get]
func (c *AmbienteController) Listar(ctx *gin.Context) {
	ambientes, err := c.servico.ListarTodos()
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar ambientes")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ambientes)
}

// BuscarPorID godoc
// @Summary      Busca um ambiente por ID
// @Description  Retorna os detalhes de um ambiente específico
// @Tags         ambiente
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Ambiente"
// @Success      200  {object}  models.Ambiente
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /ambiente/{id} [get]
func (c *AmbienteController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para buscar ambiente por ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	ambiente, err := c.servico.BuscarPorID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Ambiente não encontrado ao buscar por ID")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Ambiente não encontrado"})
		return
	}

	if err != nil {
		logrus.WithError(err).Error("Erro ao buscar ambiente por ID")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ambiente"})
		return
	}

	ctx.JSON(http.StatusOK, ambiente)
}
