package server

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/controller"
	db_infra "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/database"
	repository "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/repository"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	hostRoute       = "/api/v1"
	rotasPlantas    = "/plantas"
	rotaPlantaPorID = "/:id"
	rotaUsuarioPorID = "/:id"
)

type Server struct {
	Router *gin.Engine
}

func NewServer(db *db_infra.Database) *Server {
	router := gin.Default()

	// Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middlewares
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Repositories
	usuarioRepo := db_infra.NewUsuarioRepositorio(db.DB)
	plantaRepo := db_infra.NewPlantaRepositorio(db.DB)
	ambienteRepo := db_infra.NewAmbienteRepositorio(db.DB)
	geneticaRepo := db_infra.NewGeneticaRepositorio(db.DB)
	meioCultivoRepo := db_infra.NewMeioCultivoRepositorio(db.DB)
	diarioCultivoRepo := repository.NewDiarioCultivoRepository(db.DB)

	// Services
	usuarioService := service.NewUsuarioService(usuarioRepo)
	plantaService := service.NewPlantaService(plantaRepo, geneticaRepo, ambienteRepo, meioCultivoRepo, plantaRepo)
	ambienteService := service.NewAmbienteService(ambienteRepo)
	geneticaService := service.NewGeneticaService(geneticaRepo)
	meioCultivoService := service.NewMeioCultivoService(meioCultivoRepo)
	diarioCultivoService := service.NewDiarioCultivoService(diarioCultivoRepo, plantaRepo, ambienteRepo)

	// Controllers
	controladorUsuario := controller.NewUsuarioController(usuarioService)
	controladorPlanta := controller.NewPlantaController(plantaService)
	controladorAmbiente := controller.NewAmbienteController(ambienteService)
	controladorGenetica := controller.NewGeneticaController(geneticaService)
	controladorMeioCultivo := controller.NewMeioCultivoController(meioCultivoService)
	controladorDiarioCultivo := controller.NewDiarioCultivoController(diarioCultivoService)

	// Health check routes
	healthController := controller.NewHealthController(db.DB)
	router.GET("/health", healthController.VerificarSaude)
	router.GET("/health/ready", healthController.VerificarProntidao)
	router.GET("/health/live", healthController.VerificarVitalidade)

	// Rotas de Usuario (não autenticadas)
	router.POST(hostRoute+"/usuarios", controladorUsuario.Criar)
	router.POST(hostRoute+"/login", controladorUsuario.Login)

	// Rotas autenticadas
	authRoutes := router.Group(hostRoute)
	authRoutes.Use(middleware.AuthMiddleware())
	{
		// Rotas de Planta
		authRoutes.GET(rotasPlantas, controladorPlanta.Listar)
		authRoutes.POST(rotasPlantas, controladorPlanta.Criar)
		authRoutes.GET(rotasPlantas+rotaPlantaPorID, controladorPlanta.BuscarPorID)
		authRoutes.PUT(rotasPlantas+rotaPlantaPorID, controladorPlanta.Atualizar)
		authRoutes.DELETE(rotasPlantas+rotaPlantaPorID, controladorPlanta.Deletar)
		authRoutes.POST(rotasPlantas+rotaPlantaPorID+"/registrar-fato", controladorPlanta.RegistrarFato)

		// Rotas de Ambiente
		authRoutes.POST("/ambientes", controladorAmbiente.Criar)
		authRoutes.GET("/ambientes", controladorAmbiente.Listar)
		authRoutes.GET("/ambientes/:id", controladorAmbiente.BuscarPorID)
		authRoutes.PUT("/ambientes/:id", controladorAmbiente.Atualizar)
		authRoutes.DELETE("/ambientes/:id", controladorAmbiente.Deletar)

		// Rotas de Genetica
		authRoutes.POST("/geneticas", controladorGenetica.Criar)
		authRoutes.GET("/geneticas", controladorGenetica.Listar)
		authRoutes.GET("/geneticas/:id", controladorGenetica.BuscarPorID)
		authRoutes.PUT("/geneticas/:id", controladorGenetica.Atualizar)
		authRoutes.DELETE("/geneticas/:id", controladorGenetica.Deletar)

		// Rotas de MeioCultivo
		authRoutes.POST("/meios-cultivos", controladorMeioCultivo.Criar)
		authRoutes.GET("/meios-cultivos", controladorMeioCultivo.Listar)
		authRoutes.GET("/meios-cultivos/:id", controladorMeioCultivo.BuscarPorID)
		authRoutes.PUT("/meios-cultivos/:id", controladorMeioCultivo.Atualizar)
		authRoutes.DELETE("/meios-cultivos/:id", controladorMeioCultivo.Deletar)

		// Rotas de DiarioCultivo
		authRoutes.POST("/diarios-cultivo", controladorDiarioCultivo.Create)
		authRoutes.GET("/diarios-cultivo", controladorDiarioCultivo.List)
		authRoutes.GET("/diarios-cultivo/:id", controladorDiarioCultivo.GetByID)
		authRoutes.PUT("/diarios-cultivo/:id", controladorDiarioCultivo.Update)
		authRoutes.DELETE("/diarios-cultivo/:id", controladorDiarioCultivo.Delete)

		// Rotas de Usuario (autenticadas)
		authRoutes.GET(rotaUsuarioPorID, controladorUsuario.BuscarPorID)
		authRoutes.GET("/usuarios", controladorUsuario.Listar)
		authRoutes.PUT(rotaUsuarioPorID, controladorUsuario.Atualizar)
		authRoutes.DELETE(rotaUsuarioPorID, controladorUsuario.Deletar)
	}

	return &Server{Router: router}
}