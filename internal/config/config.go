package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
			DBHost:     getEnv("DB_HOST", "localhost"),
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
