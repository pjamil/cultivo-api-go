package controller

import (
	"errors"
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// PlantaController lida com as requisições HTTP para plantas
type PlantaController struct {
	plantaServico service.PlantaService
}

// NewPlantaController cria um novo PlantController
func NewPlantaController(plantaServico service.PlantaService) *PlantaController {
	return &PlantaController{plantaServico: plantaServico}
}

// CriarPlanta godoc
// @Summary      Cria uma nova planta
// @Description  Adiciona uma nova planta ao sistema de cultivo
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        planta  body      models.Planta  true  "Objeto da planta que precisa ser adicionado"
// @Success      201    {object}  models.Planta
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /plantas [post]
func (c *PlantaController) Criar(ctx *gin.Context) {
	var planta models.Planta
	if err := ctx.ShouldBindJSON(&planta); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Payload da requisição inválido")
		return
	}

	if err := c.plantaServico.Criar(&planta); err != nil {
		logrus.WithError(err).Error("Erro ao criar planta")
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusCreated, planta)
}

// ListarPlantas godoc
// @Summary      Lista todas as plantas
// @Description  Retorna os detalhes de todas as plantas
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Planta
// @Failure      500  {object}  map[string]interface{}
// @Router       /plantas [get]
func (c *PlantaController) Listar(ctx *gin.Context) {
	plantas, err := c.plantaServico.ListarTodas()
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar plantas")
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, plantas)
}

const (
	ErroIDPlantaInvalido             = "ID da planta inválido"
	ErroPlantaNaoEncontrada          = "Planta não encontrada"
	ErroPayloadRequisicaoInvalido    = "Payload da requisição inválido"
	ErroAtualizarPlanta              = "Erro ao atualizar a planta"
	ErroCriarPlanta                  = "Erro ao criar a planta"
	ErroRecuperarPlanta              = "Erro ao recuperar a planta"
	ErroDeletarPlanta                = "Erro ao deletar a planta"
	SucessoAtualizarPlanta           = "Planta atualizada com sucesso"
	SucessoCriarPlanta               = "Planta criada com sucesso"
	SucessoRecuperarPlanta           = "Planta recuperada com sucesso"
	SucessoListarPlantas             = "Plantas recuperadas com sucesso"
	ErroListarPlantas                = "Erro ao recuperar as plantas"
)

// BuscarPlantaPorID godoc
// @Summary      Busca uma planta por ID
// @Description  Retorna os detalhes de uma planta específica
// @ID           get-plant-by-id
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Planta"
// @Success      200  {object}  models.Planta
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /plantas/{id} [get]
func (c *PlantaController) BuscarPorID(ctx *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"id": ctx.Param("id"),
	}).Info("Buscando planta por ID")
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID da planta")
		utils.RespondWithError(ctx, http.StatusBadRequest, ErroIDPlantaInvalido)
		return
	}

	planta, err := c.plantaServico.BuscarPorID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Planta não encontrada ao buscar por ID")
		utils.RespondWithError(ctx, http.StatusNotFound, ErroPlantaNaoEncontrada)
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, planta)
}

// AtualizarPlanta godoc
// @Summary      Atualiza uma planta
// @Description  Atualiza uma planta existente
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        id      path      int            true  "ID da Planta"
// @Param        planta  body      models.Planta  true  "Objeto da planta atualizado"
// @Success      200    {object}  models.Planta
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /plantas/{id} [put]
func (c *PlantaController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID da planta para atualização")
		utils.RespondWithError(ctx, http.StatusBadRequest, ErroIDPlantaInvalido)
		return
	}
	var planta models.Planta
	if err := ctx.ShouldBindJSON(&planta); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualização de planta")
		utils.RespondWithError(ctx, http.StatusBadRequest, ErroPayloadRequisicaoInvalido)
		return
	}
	planta.ID = uint(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Planta não encontrada para atualização")
		ctx.JSON(http.StatusNotFound, gin.H{"error": ErroPlantaNaoEncontrada})
		return
	}
	if err := c.plantaServico.Atualizar(&planta); err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "update_plant",
			"plant_id":  id,
			"error":     err,
		}).Error(ErroAtualizarPlanta)
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, planta)
}

// DeletarPlanta godoc
// @Summary      Deleta uma planta
// @Description  Deleta uma planta existente
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Planta"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /plantas/{id} [delete]
func (c *PlantaController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID da planta para deleção")
		utils.RespondWithError(ctx, http.StatusBadRequest, ErroIDPlantaInvalido)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Planta não encontrada para deleção")
		ctx.JSON(http.StatusNotFound, gin.H{"error": ErroPlantaNaoEncontrada})
		return
	}
	if err := c.plantaServico.Deletar(uint(id)); err != nil {
		logrus.WithError(err).Error("Erro ao deletar planta")
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"operation": "delete_plant",
		"plant_id":  id,
	}).Info("Planta deletada com sucesso")
	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Planta deletada com sucesso"})
}
