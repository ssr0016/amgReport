package config

import (
	"amg/internal/logger"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseUrl string
	Logger      *logger.Logger
	DB          *sqlx.DB
	Pagination  PaginationConfig
	JwtSecret   string
	RedisClient *redis.Client
}

func getPort() string {
	port := os.Getenv(("HTTP_PORT"))
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("HTTP_PORT must be an int: %v\n", err)
	}

	return port
}

func getDatabaseUrl() string {
	dbUrl := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		"5432",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)

	return dbUrl
}

func getJwtSecret() string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET is not set in the environment")
	}
	return jwtSecret
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	// Initialize Logger
	development := os.Getenv("ENVIRONMENT") == "development"

	loggers, err := logger.Init(development)
	if err != nil {
		log.Fatalf("Error initializing logger: %v\n", err)
	}

	// Initialize Database
	dbConn, err := sqlx.Connect("postgres", getDatabaseUrl())
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	// Initialize Redis Client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Verify Redis connection
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v\n", err)
	}

	cfg := &Config{
		Port:        getPort(),
		DatabaseUrl: getDatabaseUrl(),
		Logger:      loggers,
		DB:          dbConn,
		JwtSecret:   getJwtSecret(),
		RedisClient: redisClient,
	}
	cfg.Logger = loggers

	// Apply pagination config
	cfg.LoadPaginationConfig()

	return cfg
}
