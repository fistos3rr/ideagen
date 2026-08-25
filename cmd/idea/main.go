package main

import (
	"flag"
	"os"
	"errors"

	"github.com/fistos3rr/ideagen/internal/jsonlog"
	"github.com/fistos3rr/ideagen/internal/ai"
)

type config struct {
	port int
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

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development",
		"Environment (development|staging|production)")
	flag.StringVar(&cfg.aiProviderType, "ai", "groq", "AI chat provider")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	var provider ai.Provider
	switch cfg.aiProviderType {
	case "groq":
		aicfg := ai.Config{
			APIKey: os.Getenv("API_KEY_GROQ_IDEAGEN"),
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
