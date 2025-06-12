package server

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/controller"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
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
	plantController := controller.NewPlantaController(plantService)

	const hostRoute = "/api/v1"
	const plantsRoute = "/api/v1/plants"
	const plantByIDRoute = "/:id"
	router.GET(plantsRoute, plantController.GetAllPlants)
	router.POST(plantsRoute, plantController.CreatePlanta)
	router.GET(plantsRoute+plantByIDRoute, plantController.GetPlantByID)
	router.PUT(plantsRoute+plantByIDRoute, plantController.UpdatePlant)
	router.DELETE(plantsRoute+plantByIDRoute, plantController.DeletePlant)

	// Ambiente routes
	router.POST(hostRoute+"/ambientes", controller.CreateAmbiente(db.DB))
	router.GET(hostRoute+"/ambientes", controller.ListAmbientes(db.DB))

	// Genetica routes
	geneticaRepo := repository.NewGeneticaRepository(db.DB)
	geneticaService := service.NewGeneticaService(geneticaRepo)
	geneticaController := controller.NewGeneticaController(geneticaService)
	router.POST(hostRoute+"/geneticas", geneticaController.Create)
	router.GET(hostRoute+"/geneticas", geneticaController.GetAll)

	meioCultivoRepo := repository.NewMeioCultivoRepository(db.DB)
	meioCultivoService := service.NewMeioCultivoService(meioCultivoRepo)
	meioCultivoController := controller.NewMeioCultivoController(meioCultivoService)
	router.POST(hostRoute+"/meios_cultivo", meioCultivoController.Create)

	return &Server{Router: router}
}
