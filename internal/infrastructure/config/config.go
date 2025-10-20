package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"getswing.app/player-service/internal/constants"
)

type Config struct {
	HTTPPort     string
	DatabaseURL  string
	DBName       string
	DBUser       string
	DBPassword   string
	DBHost       string
	DBSslMode    string
	DBPort       int
	RabbitURL    string
	LogLevel     string
	JWTSecret    string
	JWTIssuer    string
	JWTAccessTTL time.Duration
	GrpcServer   string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func getenvDuration(key, def string) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	if d, err := time.ParseDuration(def); err == nil {
		return d
	}
	// fallback 15m
	return 15 * time.Minute
}

func Load() Config {

	LoadEnvFile(".env")

	return Config{
		HTTPPort:     getenv(constants.EnvHTTPPort, "3000"),
		DBName:       getenv(constants.EnvDBName, "player_service"),
		DBUser:       getenv(constants.EnvDBUser, "swing"),
		DBPassword:   getenv(constants.EnvDBPass, "password"),
		DBHost:       getenv(constants.EnvDBHost, "localhost"),
		DBPort:       getenvInt(constants.EnvDBPort, 5432),
		RabbitURL:    getenv(constants.EnvRabbitURL, "amqp://guest:guest@localhost:5672/"),
		LogLevel:     getenv(constants.EnvLogLevel, "info"),
		JWTSecret:    getenv(constants.EnvJWTSecret, "devsecret-change-me"),
		JWTIssuer:    getenv(constants.EnvJWTIssuer, "player-service"),
		JWTAccessTTL: getenvDuration(constants.EnvJWTAccessTTL, "15m"),
		GrpcServer:   getenv(constants.EnvGrpcServer, ""),
	}
}

func LoadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("⚠️  No %s file found, skipping .env loading", filename)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments or empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Split key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"")

		// Set environment variable
		os.Setenv(key, val)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}
}
