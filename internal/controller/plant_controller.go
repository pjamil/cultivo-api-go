package controller

import (
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
)

// PlantController handles HTTP requests for plants
type PlantController struct {
	plantService *service.PlantService
}

// NewPlantController creates a new PlantController
func NewPlantController(plantService *service.PlantService) *PlantController {
	return &PlantController{plantService: plantService}
}

// CreatePlant godoc
// @Summary Create a new plant
// @Description Add a new plant to the cultivation system
// @Tags plants
// @Accept  json
// @Produce  json
// @Param plant body models.Plant true "Plant object that needs to be added"
// @Success 201 {object} models.Plant
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /plants [post]
func (c *PlantController) CreatePlant(ctx *gin.Context) {
	var plant models.Plant
	if err := ctx.ShouldBindJSON(&plant); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := c.plantService.CreatePlant(&plant); err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusCreated, plant)
}

// GetAllPlants godoc
// @Summary Get all plants
// @Description Get details of all plants
// @Tags plants
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Plant
// @Failure 500 {object} map[string]interface{}
// @Router /plants [get]
func (c *PlantController) GetAllPlants(ctx *gin.Context) {
	plants, err := c.plantService.GetAllPlants()
	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, plants)
}

// GetPlantByID godoc
// @Summary Get plant by ID
// @Description Get details of a specific plant
// @Tags plants
// @Accept  json
// @Produce  json
// @Param id path int true "Plant ID"
// @Success 200 {object} models.Plant
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /plants/{id} [get]
func (c *PlantController) GetPlantByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	plant, err := c.plantService.GetPlantByID(uint(id))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusNotFound, "Plant not found")
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
// @Param id path int true "Plant ID"
// @Param plant body models.Plant true "Updated plant object"
// @Success 200 {object} models.Plant
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /plants/{id} [put]
func (c *PlantController) UpdatePlant(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	var plant models.Plant
	if err := ctx.ShouldBindJSON(&plant); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	plant.ID = uint(id)
	if err := c.plantService.UpdatePlant(&plant); err != nil {
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
// @Param id path int true "Plant ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /plants/{id} [delete]
func (c *PlantController) DeletePlant(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	if err := c.plantService.DeletePlant(uint(id)); err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Plant deleted successfully"})
}
