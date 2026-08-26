package main

import (
	"os"
	"errors"
	"strconv"
	"database/sql"
	"time"
	"context"

	"github.com/fistos3rr/ideagen/internal/jsonlog"
	"github.com/fistos3rr/ideagen/internal/ai"
	"github.com/fistos3rr/ideagen/internal/data"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	env string
	aiProviderType string
	db struct {
		dsn string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	aiProvider ai.Provider
	models data.Models
}

func (cfg *Config) parseEnv() {
	strVal := os.Getenv("PORT")
	cfg.port, err := strconv.Atoi(strVal)
	if err != nil {
		port = 8080
	}

	strVal = os.Getenv("AI_PROVIDER")
	cfg.aiProviderType = strVal
	if cfg.aiProviderType == "" {
		cfg.aiProviderType = "groq"
	}

	strVal = os.Getenv("DB_DSN")
	cfg.db.dsn = strVal

	strVal = os.Getenv("DB_MAX_OPEN_CONNS")
	cfg.db.maxOpenConns, err := strconv.Atoi(strVal)
	if err != nil {
		cfg.db.maxOpenConns = 25
	}

	strVal = os.Getenv("DB_MAX_IDLE_CONNS")
	cfg.db.maxIdleConns, err = strconv.Atoi(strVal)
	if err != nil {
		cfg.db.maxIdleConns = 25
	}

	strVal = os.Getenv("DB_MAX_IDLE_TIME")
	cfg.db.maxIdleTime = strVal
	if cfg.db.maxIdleTime == "" {
		cfg.db.maxIdleTime = "15m"
	}
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

	var provider ai.Provider
	switch cfg.aiProviderType {
	case "groq":
		aicfg := ai.Config{
			APIKey: os.Getenv("AI_API_KEY"),
			Model: "openai/gpt-oss-20b",
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
		config: cfg,
		logger: logger,
		aiProvider: provider,
		models: data.NewModels(db),
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
