package server

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/controller"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/database"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Router *gin.Engine
}

func NewServer(db *database.Database) *Server {
	router := gin.Default()

	// Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middlewares
	router.Use(middleware.LoggingMiddleware())

	// Health check routes
	healthController := controller.NewHealthController()
	router.GET("/health", healthController.CheckHealth)
	router.GET("/health/ready", healthController.ReadyCheck)
	router.GET("/health/live", healthController.LiveCheck)

	// Plant routes
	plantRepo := database.NewPlantRepository(db.DB)
	plantService := service.NewPlantService(plantRepo)
	plantController := controller.NewPlantController(plantService)

	const plantByIDRoute = "/plants/:id"
	api := router.Group("/api/v1")
	{
		api.GET("/plants", plantController.GetAllPlants)
		api.POST("/plants", plantController.CreatePlant)
		api.GET(plantByIDRoute, plantController.GetPlantByID)
		api.PUT(plantByIDRoute, plantController.UpdatePlant)
		api.DELETE(plantByIDRoute, plantController.DeletePlant)
	}

	return &Server{Router: router}
}
