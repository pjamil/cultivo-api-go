package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	// Database configuration
	// These values can be set in a .env file or as environment variables
	// Default values are provided for local development
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		log.Println("Using environment variables from the system")
	}

	appEnv := getEnv("APP_ENV", "development")

	config := &Config{}
	log.Println(config)
	if appEnv == "test" {
		config = &Config{
			DBDriver:   getEnv("DB_DRIVER_TEST", "postgres"),
			DBHost:     getEnv("DB_HOST_TEST", "db_test"),
			DBPort:     getEnv("DB_PORT_TEST", "5432"), // A porta do container é 5432
			DBUser:     getEnv("DB_USER_TEST", "testuser"),
			DBPassword: getEnv("DB_PASSWORD_TEST", "testpassword"),
			DBName:     getEnv("DB_NAME_TEST", "cultivo_test_db"),
			ServerPort: getEnv("SERVER_PORT", "8080"),
		}
	} else {
		config = &Config{
			DBDriver:   getEnv("DB_DRIVER", "postgres"),
			DBHost:     getEnv("DB_HOST", "db"),
			DBPort:     getEnv("DB_PORT", "5435"),
			DBUser:     getEnv("DB_USER", "paulo"),
			DBPassword: getEnv("DB_PASSWORD", "753951465827PJamil"),
			DBName:     getEnv("DB_NAME", "cultivo-api-go"),
			ServerPort: getEnv("SERVER_PORT", "8080"),
		}
	}

	log.Println(config)
	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// ConnectDatabase inicializa e retorna uma conexão GORM com o banco de dados.
func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	log.Println("Conexão com o banco de dados estabelecida")
	return db, nil
}
