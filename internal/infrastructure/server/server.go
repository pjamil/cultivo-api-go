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
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Health check routes
	healthController := controller.NewHealthController()
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

		// Rotas de Usuario (autenticadas)
		authRoutes.GET(rotaUsuarioPorID, controladorUsuario.BuscarPorID)
		authRoutes.GET("/usuarios", controladorUsuario.Listar)
		authRoutes.PUT(rotaUsuarioPorID, controladorUsuario.Atualizar)
		authRoutes.DELETE(rotaUsuarioPorID, controladorUsuario.Deletar)
	}

	return &Server{Router: router}
}
