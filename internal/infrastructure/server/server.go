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
	geneticaRepo := repository.NewGeneticaRepository(db.DB)
	ambienteRepo := repository.NewAmbienteRepository(db.DB)
	meioCultivoRepo := repository.NewMeioCultivoRepository(db.DB)
	plantService := service.NewPlantService(plantRepo, geneticaRepo, ambienteRepo, meioCultivoRepo)
	plantController := controller.NewPlantaController(plantService)

	const hostRoute = "/api/v1"
	const plantsRoute = "/api/v1/plants"
	const plantByIDRoute = "/:id"
	const usuarioByIDRoute = hostRoute + "/usuarios/:id"
	router.GET(plantsRoute, plantController.GetAllPlants)
	router.POST(plantsRoute, plantController.CreatePlanta)
	router.GET(plantsRoute+plantByIDRoute, plantController.GetPlantByID)
	router.PUT(plantsRoute+plantByIDRoute, plantController.UpdatePlant)
	router.DELETE(plantsRoute+plantByIDRoute, plantController.DeletePlant)

	// Ambiente routes
	ambienteRepo = repository.NewAmbienteRepository(db.DB)
	ambienteService := service.NewAmbienteService(ambienteRepo)
	ambienteController := controller.NewAmbienteController(ambienteService)
	router.POST(hostRoute+"/ambientes", ambienteController.CreateAmbiente)
	router.GET(hostRoute+"/ambientes", ambienteController.ListAmbientes)
	router.GET(hostRoute+"/ambientes/:id", ambienteController.GetAmbienteByID)

	// Genetica routes
	geneticaRepo = repository.NewGeneticaRepository(db.DB)
	geneticaService := service.NewGeneticaService(geneticaRepo)
	geneticaController := controller.NewGeneticaController(geneticaService)
	router.POST(hostRoute+"/geneticas", geneticaController.Create)
	router.GET(hostRoute+"/geneticas", geneticaController.GetAll)
	router.GET(hostRoute+"/geneticas/:id", geneticaController.GetGeneticaByID)

	// MeioCultivo routes
	meioCultivoRepo = repository.NewMeioCultivoRepository(db.DB)
	meioCultivoService := service.NewMeioCultivoService(meioCultivoRepo)
	meioCultivoController := controller.NewMeioCultivoController(meioCultivoService)
	router.POST(hostRoute+"/meios_cultivo", meioCultivoController.Create)
	router.GET(hostRoute+"/meios_cultivo", meioCultivoController.GetAll)
	router.GET(hostRoute+"/meios_cultivo/:id", meioCultivoController.GetByID)

	// Usuario routes
	usuarioRepo := repository.NewUsuarioRepository(db.DB)
	usuarioService := service.NewUsuarioService(usuarioRepo)
	usuarioController := controller.NewUsuarioController(usuarioService)
	router.POST(hostRoute+"/usuarios", usuarioController.Create)
	router.GET(usuarioByIDRoute, usuarioController.GetByID)
	router.GET(hostRoute+"/usuarios", usuarioController.GetAll)
	router.PUT(usuarioByIDRoute, usuarioController.Update)
	router.DELETE(usuarioByIDRoute, usuarioController.Delete)

	return &Server{Router: router}
}
