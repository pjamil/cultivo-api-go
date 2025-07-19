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

	"github.com/joho/godotenv"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/config"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/database"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/server"
	"github.com/sirupsen/logrus"
)

// @title Cultivo API
// @version 1.0
// @description This is the API for the Cultivo app.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	cfg := config.LoadConfig()

	db, err := database.NewDatabase(cfg)
	if err != nil {
		logrus.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	// Setup do Router
	servidor := server.NewServer(db)

	// Graceful Shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: servidor.Router,
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
