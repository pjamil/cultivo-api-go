package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/config"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/repository"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/handler"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/router"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	cfg := config.LoadConfig()
	

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		logrus.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	

	// =========================================================================
	// Injeção de Dependências (DI)
	// =========================================================================
	// Construímos nossas dependências da camada mais interna para a mais externa:
	// Database -> Repositório -> Serviço -> Handler
	diarioRepo := repository.NewDiarioCultivoRepository(db)
	diarioService := service.NewDiarioCultivoService(diarioRepo)
	diarioHandler := handler.NewDiarioCultivoHandler(diarioService)

	// Setup do Router
	r := gin.New()
	r.Use(gin.Recovery())
	// Usando o middleware de log customizado que já formata para Logrus
	r.Use(middleware.LoggerMiddleware())

	// Rota pública de health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Grupo de rotas da API v1
	apiV1 := r.Group("/api/v1")
	{
		// Rotas públicas dentro da API
		// TODO: Adicionar rotas de autenticação (login, registro) aqui

		// Grupo de rotas que exigem autenticação
		authorized := apiV1.Group("/")
		
		{
			// As rotas de diário de cultivo são registradas aqui, dentro do grupo protegido.
			router.RegisterDiarioCultivoRoutes(authorized, diarioHandler)

			// Futuras rotas protegidas (ex: plantas, ambientes) serão registradas aqui.
		}
	}

	// Graceful Shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: r,
	}

	go func() {
		logrus.Infof("Servidor iniciado na porta %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Println("Desligando o servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Erro no desligamento do servidor:", err)
	}

	logrus.Println("Servidor desligado com sucesso.")
}
