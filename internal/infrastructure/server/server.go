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
	router.GET("/health", healthController.VerificarSaude)
	router.GET("/health/ready", healthController.VerificarProntidao)
	router.GET("/health/live", healthController.VerificarVitalidade)

	// Rotas de Planta
	repositorioPlanta := database.NewPlantaRepositorio(db.DB)
	geneticaRepo := repository.NewGeneticaRepositorio(db.DB)
	ambienteRepo := repository.NewAmbienteRepositorio(db.DB)
	meioCultivoRepo := repository.NewMeioCultivoRepositorio(db.DB)
	servicoPlanta := service.NewPlantaService(repositorioPlanta, geneticaRepo, ambienteRepo, meioCultivoRepo)
	controladorPlanta := controller.NewPlantaController(servicoPlanta)

	const hostRoute = "/api/v1"
	const rotasPlantas = "/api/v1/plantas"
	const rotaPlantaPorID = "/:id"
	const rotaUsuarioPorID = hostRoute + "/usuarios/:id"
	router.GET(rotasPlantas, controladorPlanta.Listar)
	router.POST(rotasPlantas, controladorPlanta.Criar)
	router.GET(rotasPlantas+rotaPlantaPorID, controladorPlanta.BuscarPorID)
	router.PUT(rotasPlantas+rotaPlantaPorID, controladorPlanta.Atualizar)
	router.DELETE(rotasPlantas+rotaPlantaPorID, controladorPlanta.Deletar)

	// Rotas de Ambiente
	repositorioAmbiente := repository.NewAmbienteRepositorio(db.DB)
	servicoAmbiente := service.NewAmbienteService(repositorioAmbiente)
	controladorAmbiente := controller.NewAmbienteController(servicoAmbiente)
	router.POST(hostRoute+"/ambientes", controladorAmbiente.Criar)
	router.GET(hostRoute+"/ambientes", controladorAmbiente.Listar)
	router.GET(hostRoute+"/ambientes/:id", controladorAmbiente.BuscarPorID)
	router.PUT(hostRoute+"/ambientes/:id", controladorAmbiente.Atualizar)
	router.DELETE(hostRoute+"/ambientes/:id", controladorAmbiente.Deletar)

	// Rotas de Genetica
	repositorioGenetica := repository.NewGeneticaRepositorio(db.DB)
	servicoGenetica := service.NewGeneticaService(repositorioGenetica)
	controladorGenetica := controller.NewGeneticaController(servicoGenetica)
	router.POST(hostRoute+"/geneticas", controladorGenetica.Criar)
	router.GET(hostRoute+"/geneticas", controladorGenetica.Listar)
	router.GET(hostRoute+"/geneticas/:id", controladorGenetica.BuscarPorID)
	router.PUT(hostRoute+"/geneticas/:id", controladorGenetica.Atualizar)
	router.DELETE(hostRoute+"/geneticas/:id", controladorGenetica.Deletar)

	// Rotas de MeioCultivo
	repositorioMeioCultivo := repository.NewMeioCultivoRepositorio(db.DB)
	servicoMeioCultivo := service.NewMeioCultivoService(repositorioMeioCultivo)
	controladorMeioCultivo := controller.NewMeioCultivoController(servicoMeioCultivo)
	router.POST(hostRoute+"/meios-cultivos", controladorMeioCultivo.Criar)
	router.GET(hostRoute+"/meios-cultivos", controladorMeioCultivo.Listar)
	router.GET(hostRoute+"/meios-cultivos/:id", controladorMeioCultivo.BuscarPorID)
	router.PUT(hostRoute+"/meios-cultivos/:id", controladorMeioCultivo.Atualizar)
	router.DELETE(hostRoute+"/meios-cultivos/:id", controladorMeioCultivo.Deletar)

	// Rotas de Usuario
	repositorioUsuario := repository.NewUsuarioRepositorio(db.DB)
	servicoUsuario := service.NewUsuarioService(repositorioUsuario)
	controladorUsuario := controller.NewUsuarioController(servicoUsuario)
	router.POST(hostRoute+"/usuarios", controladorUsuario.Criar)
	router.GET(rotaUsuarioPorID, controladorUsuario.BuscarPorID)
	router.GET(hostRoute+"/usuarios", controladorUsuario.Listar)
	router.PUT(rotaUsuarioPorID, controladorUsuario.Atualizar)
	router.DELETE(rotaUsuarioPorID, controladorUsuario.Deletar)

	return &Server{Router: router}
}
