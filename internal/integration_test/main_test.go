package integration_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/config"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/database"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/infrastructure/server"
	"github.com/gin-gonic/gin"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var testServer *server.Server
var testDB *database.Database

func TestMain(m *testing.M) {
	// Definir o ambiente para "test"
	os.Setenv("APP_ENV", "test")
	os.Setenv("DB_HOST_TEST", "localhost") // Usar localhost para testes locais com docker-compose
	os.Setenv("DB_PORT_TEST", "5435")     // Porta do docker-compose.test.yml
	os.Setenv("DB_USER_TEST", "testuser")
	os.Setenv("DB_PASSWORD_TEST", "testpassword")
	os.Setenv("DB_NAME_TEST", "cultivo_test_db")

	// Carregar a configuração de teste
	cfg := config.LoadConfig()

	// Inicializar o banco de dados de teste
	var err error
	testDB, err = database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Executar migrações no banco de dados de teste
	log.Println("Running migrations for test database...")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	m, err := migrate.New(
		"file://../../infrastructure/database/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied successfully.")

	// Criar o servidor Gin para testes
	gin.SetMode(gin.TestMode) // Define o modo Gin para teste
	testServer = server.NewServer(testDB)

	// Rodar os testes
	code := m.Run()

	// Limpar o banco de dados de teste após os testes (opcional, mas recomendado)
	log.Println("Cleaning up test database...")
	// Para garantir que o banco esteja limpo para a próxima execução, podemos dar um Down e Up
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Failed to clean up migrations: %v", err)
	}
	testDB.DB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;") // Limpeza simples

	os.Exit(code)
}

// GetTestRouter retorna o router do servidor de teste
func GetTestRouter() *gin.Engine {
	return testServer.Router
}

// GetTestDB retorna a instância do banco de dados de teste
func GetTestDB() *database.Database {
	return testDB
}
