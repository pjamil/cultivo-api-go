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

// PlantaController handles HTTP requests for plants
type PlantaController struct {
	plantaService *service.PlantaService
}

// NewPlantaController creates a new PlantController
func NewPlantaController(plantaService *service.PlantaService) *PlantaController {
	return &PlantaController{plantaService: plantaService}
}

// CreatePlanta godoc
// @Summary Create a new plant
// @Description Add a new plant to the cultivation system
// @Tags plants
// @Accept  json
// @Produce  json
// @Param plant body models.Planta true "Planta object that needs to be added"
// @Success 201 {object} models.Planta
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /plants [post]
func (c *PlantaController) CreatePlanta(ctx *gin.Context) {
	var planta models.Planta
	if err := ctx.ShouldBindJSON(&planta); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := c.plantaService.CreatePlanta(&planta); err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusCreated, planta)
}

// GetAllPlants godoc
// @Summary Get all plants
// @Description Get details of all plants
// @Tags plants
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Planta
// @Failure 500 {object} map[string]interface{}
// @Router /plants [get]
func (c *PlantaController) GetAllPlants(ctx *gin.Context) {
	plants, err := c.plantaService.GetAllPlants()
	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, plants)
}

const (
	InvalidPlantIDError              = "Invalid plant ID"
	PlantNotFoundError               = "Planta not found"
	InvalidRequestPayloadError       = "Invalid request payload"
	PlantUpdateError                 = "Error updating plant"
	PlantCreationError               = "Error creating plant"
	PlantRetrievalError              = "Error retrieving plant"
	PlantDeletionError               = "Error deleting plant"
	PlantUpdateSuccessMessage        = "Planta updated successfully"
	PlantCreationSuccessMessage      = "Planta created successfully"
	PlantRetrievalSuccessMessage     = "Planta retrieved successfully"
	PlantListRetrievalSuccessMessage = "Plants retrieved successfully"
	PlantListRetrievalError          = "Error retrieving plants"
	PlantCreationSuccess             = "Planta created successfully"
	PlantUpdateSuccess               = "Planta updated successfully"
	PlantDeletionSuccess             = "Planta deleted successfully"
	PlantRetrievalSuccess            = "Planta retrieved successfully"
	PlantListRetrievalSuccess        = "Plants retrieved successfully"
)

// GetPlantByID godoc
// @Summary Get plant by ID
// @Description Get details of a specific plant
// @ID get-plant-by-id
// @Tags plants
// @Accept  json
// @Produce  json
// @Param id path int true "Planta ID"
// @Success 200 {object} models.Planta
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /plants/{id} [get]
func (c *PlantaController) GetPlantByID(ctx *gin.Context) {
	// Log the request for fetching a plant by ID
	logrus.WithFields(logrus.Fields{
		"id": ctx.Param("id"),
	}).Info("Fetching plant")
	// Parse the plant ID from the URL parameter
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	plant, err := c.plantaService.GetPlantByID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.RespondWithError(ctx, http.StatusNotFound, "Planta não encontrada")
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, plant)
}

// UpdatePlant godoc
// @Summary Update a plant
// @Description Update an existing plant
// @Tags plants
// @Accept  json
// @Produce  json
// @Param id path int true "Planta ID"
// @Param plant body models.Planta true "Updated plant object"
// @Success 200 {object} models.Planta
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /plants/{id} [put]
func (c *PlantaController) UpdatePlant(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, InvalidPlantIDError)
		return
	}
	var plant models.Planta
	if err := ctx.ShouldBindJSON(&plant); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, InvalidRequestPayloadError)
		return
	}
	plant.ID = uint(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Planta não encontrada"})
		return
	}
	if err := c.plantaService.UpdatePlant(&plant); err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "update_plant",
			"plant_id":  id,
			"error":     err,
		}).Error("Erro ao atualizar planta")
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, plant)
}

// DeletePlant godoc
// @Summary Delete a plant
// @Description Delete an existing plant
// @Tags plants
// @Accept  json
// @Produce  json
// @Param id path int true "Planta ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /plants/{id} [delete]
func (c *PlantaController) DeletePlant(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, InvalidPlantIDError)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Planta não encontrada"})
		return
	}
	if err := c.plantaService.DeletePlant(uint(id)); err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"operation": "delete_plant",
		"plant_id":  id,
	}).Info("Planta deletada com sucesso")
	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Planta deleted successfully"})
}
