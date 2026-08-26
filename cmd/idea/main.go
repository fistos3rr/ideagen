package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fistos3rr/ideagen/internal/ai"
	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/jsonlog"

	_ "github.com/lib/pq"
)

type config struct {
	port           int
	env            string
	aiProviderType string
	db             struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config     config
	logger     *jsonlog.Logger
	aiProvider ai.Provider
	models     data.Models
	wg sync.WaitGroup
}

func (cfg *config) parseEnv() {
	strVal := os.Getenv("PORT")
	port, err := strconv.Atoi(strVal)
	if err != nil {
		port = 8080
	}

	strVal = os.Getenv("AI_PROVIDER")
	aiProviderType := strVal
	if aiProviderType == "" {
		aiProviderType = "groq"
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	strVal = os.Getenv("DB_MAX_OPEN_CONNS")
	maxOpenConns, err := strconv.Atoi(strVal)
	if err != nil {
		maxOpenConns = 25
	}

	strVal = os.Getenv("DB_MAX_IDLE_CONNS")
	maxIdleConns, err := strconv.Atoi(strVal)
	if err != nil {
		maxIdleConns = 25
	}

	strVal = os.Getenv("DB_MAX_IDLE_TIME")
	maxIdleTime := strVal
	if maxIdleTime == "" {
		maxIdleTime = "15m"
	}

	cfg.port = port
	cfg.aiProviderType = aiProviderType
	cfg.db.dsn = dsn
	cfg.db.maxOpenConns = maxOpenConns
	cfg.db.maxIdleConns = maxIdleConns
	cfg.db.maxIdleTime = maxIdleTime
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	var cfg config
	cfg.parseEnv()
	cfg.env = "development"

	logger.PrintInfo("database config", map[string]string{
		"dsn": cfg.db.dsn,
	})

	var provider ai.Provider
	switch cfg.aiProviderType {
	case "groq":
		aicfg := ai.Config{
			APIKey: os.Getenv("AI_API_KEY"),
			Model:  "openai/gpt-oss-20b",
		}

		logData := map[string]string{

			"model": aicfg.Model,
		}

		if cfg.env == "development" {
			logData["api_key"] = aicfg.APIKey
		}

		logger.PrintInfo("running groq ai provider", logData)
		provider = ai.NewGroqClient(aicfg)
	default:
		logger.PrintFatal(errors.New("unknown provider"), map[string]string{
			"provider": cfg.aiProviderType,
		})
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config:     cfg,
		logger:     logger,
		aiProvider: provider,
		models:     data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
