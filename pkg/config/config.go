package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Kafka
	KafkaBootstrapServers string
	KafkaGroupID          string

	// gRPC
	GRPCPort int
}

func LoadConfig() *Config {
	return &Config{
		// Database defaults
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "trendyol_user"),
		DBPassword: getEnv("DB_PASSWORD", "trendyol_password"),
		DBName:     getEnv("DB_NAME", "trendyol_tracker"),

		// Kafka defaults
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		KafkaGroupID:          getEnv("KAFKA_GROUP_ID", "trendyol-tracker"),

		// gRPC defaults
		GRPCPort: getEnvAsInt("GRPC_PORT", 50051),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
