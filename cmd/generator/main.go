package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fistos3rr/ideagen/internal/ai"
	"github.com/fistos3rr/ideagen/internal/jsonlog"
	"github.com/fistos3rr/ideagen/internal/prompt"
	"github.com/spf13/viper"
)

var logger *jsonlog.Logger

func getPrompts(t string) (string, string) {
	promptManager, err := prompt.NewPromptManager(".")
	ok := prompt.IsDefaultErr(err)
	if ok {
		logger.PrintInfo(err.Error(), nil)
	} else if err != nil {
		panic(err)
	}

	sysPr, pr, err := promptManager.GetPrompts(t)
	if err != nil {
		panic(err)
	}

	return sysPr, pr
}

func main() {
	logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)	

	providerType := flag.String("provider", "groq", "ai provider type")
	t := flag.String("type", "", "topic type")
	flag.Parse()

	if len(*t) == 0 {
		panic("Provide topic type")
	}

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading app.env file: %w", err))
	}

	aicfg := ai.Config{
		APIKey: viper.GetString("AI_API_KEY"),
		Model:  viper.GetString("AI_MODEL"),
		APIURL: viper.GetString("AI_API_URL"),
	}

	sysPr, pr := getPrompts(*t)

	var provider ai.Provider
	switch *providerType {
	case "groq":
		request := ai.NewGroqRequest(aicfg)
		provider = ai.NewGroqClientWithRequest(aicfg, request)
		if len(sysPr) > 0 {
			request.AddSystemMessage(sysPr)
		}
		request.AddMessage(pr)
	default:
		panic("unknown provider type")
	}

	answer, err := provider.SendRequest(context.Background())
	if err != nil {
		panic(fmt.Errorf("error while asking ai: %w", err))
	}

	fmt.Println("PROMPT:")
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println(pr)
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println()
	fmt.Println("ANSWER:")
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println(answer)
	fmt.Println(strings.Repeat("=", 20))
}
