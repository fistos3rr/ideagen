package main

import (
	"os"
	"errors"

	"github.com/fistos3rr/ideagen/internal/jsonlog"
	"github.com/fistos3rr/ideagen/internal/ai"
)

type config struct {
	port string
	env string
	aiProviderType string
}

type application struct {
	config config
	logger *jsonlog.Logger
	aiProvider ai.Provider
}

func main() {
	var cfg config


	cfg.port = os.Getenv("PORT")
	cfg.env = "development"
	cfg.aiProviderType = os.Getenv("AI_PROVIDER")

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

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

	app := &application{
		config: cfg,
		logger: logger,
		aiProvider: provider,
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
