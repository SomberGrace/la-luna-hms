// backend/internal/config/config.go
package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime configuration for the La Luna HMS server.
type Config struct {
	Port                 string
	JWTSecret            string
	AllowDemoUsers       bool
	DatabaseURL          string
	DBMaxOpenConns       int
	DBMaxIdleConns       int
	DBConnMaxLifetime    time.Duration
	DBConnMaxIdleTime    time.Duration
	RunStartupMigrations bool
	SeedDatabase         bool
}

// LoadConfig reads environment variables (with .env fallback) and returns a Config.
func LoadConfig() (*Config, error) {
	loadEnvFile(".env")
	loadEnvFile("../../.env")
	loadEnvFile("../../../.env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	allowDemoUsers := strings.EqualFold(strings.TrimSpace(os.Getenv("ALLOW_DEMO_USERS")), "true")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("SUPABASE_DB_URL")
	}

	dbMaxOpenConns, err := intFromEnv("DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		return nil, err
	}
	if dbMaxOpenConns <= 0 {
		return nil, fmt.Errorf("DB_MAX_OPEN_CONNS must be greater than zero")
	}

	dbMaxIdleConns, err := intFromEnv("DB_MAX_IDLE_CONNS", 2)
	if err != nil {
		return nil, err
	}
	if dbMaxIdleConns < 0 || dbMaxIdleConns > dbMaxOpenConns {
		return nil, fmt.Errorf("DB_MAX_IDLE_CONNS must be between zero and DB_MAX_OPEN_CONNS")
	}

	dbConnMaxLifetime, err := durationFromEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute)
	if err != nil {
		return nil, err
	}
	if dbConnMaxLifetime < 0 {
		return nil, fmt.Errorf("DB_CONN_MAX_LIFETIME cannot be negative")
	}

	dbConnMaxIdleTime, err := durationFromEnv("DB_CONN_MAX_IDLE_TIME", time.Minute)
	if err != nil {
		return nil, err
	}
	if dbConnMaxIdleTime < 0 {
		return nil, fmt.Errorf("DB_CONN_MAX_IDLE_TIME cannot be negative")
	}

	runStartupMigrations, err := boolFromEnv("RUN_STARTUP_MIGRATIONS", true)
	if err != nil {
		return nil, err
	}

	seedDatabase, err := boolFromEnv("SEED_DATABASE", true)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:                 port,
		JWTSecret:            jwtSecret,
		AllowDemoUsers:       allowDemoUsers,
		DatabaseURL:          dbURL,
		DBMaxOpenConns:       dbMaxOpenConns,
		DBMaxIdleConns:       dbMaxIdleConns,
		DBConnMaxLifetime:    dbConnMaxLifetime,
		DBConnMaxIdleTime:    dbConnMaxIdleTime,
		RunStartupMigrations: runStartupMigrations,
		SeedDatabase:         seedDatabase,
	}, nil
}

func intFromEnv(name string, defaultValue int) (int, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", name, err)
	}
	return parsed, nil
}

func durationFromEnv(name string, defaultValue time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", name, err)
	}
	return parsed, nil
}

func boolFromEnv(name string, defaultValue bool) (bool, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", name, err)
	}
	return parsed, nil
}

// loadEnvFile reads a KEY=VALUE .env file and sets environment variables
// for any keys not already set in the environment.
func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}
