package main

import (
	"net/http"
	"os"

	_ "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/docs" // Import the generated Swagger docs
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/config"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/database"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/server"

	"github.com/sirupsen/logrus"
)

// @title Plant Cultivation API
// @version 1.0
// @description API para gerenciamento de cultivo de plantas
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@plantcultivation.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

func main() {
	// Configure Logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	// Create server
	srv := server.NewServer(db)

	// Start server
	logrus.Printf("Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, srv.Router); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
