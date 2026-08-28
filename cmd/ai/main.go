package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fistos3rr/ideagen/internal/ai"
	"github.com/spf13/viper"
)

type fileString string

func (fs *fileString) String() string {
	return string(*fs)
}

func (fs *fileString) Set(value string) error {
	if strings.HasPrefix(value, "./") {
		content, err := os.ReadFile(value[2:])
		if err != nil {
			return err
		}
		*fs = fileString(content)
	} else {
		*fs = fileString(value)
	}
	return nil
}

func main() {
	data := flag.String("prompt", "./prompt.md", "prompt (if starts with ./ will use file instead of stdin text)")
	dataSys := flag.String("system", "./system.md", "system prompt")
	providerType := flag.String("provider", "groq", "ai provider type")
	flag.Parse()

	if len(*data) == 0 {
		fmt.Println("Please provide prompt.")
		os.Exit(1)
	}

	var prompt fileString
	err := prompt.Set(*data)
	if err != nil {
		panic(fmt.Errorf("error reading prompt: %w\n", err))
	}

	var systemPrompt fileString
	err = systemPrompt.Set(*dataSys)
	if err != nil {
		panic(fmt.Errorf("error reading prompt: %w\n", err))
	}

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading app.env file: %w\n", err))
	}

	aicfg := ai.Config{
		APIKey: viper.GetString("AI_API_KEY"),
		Model:  viper.GetString("AI_MODEL"),
		APIURL: viper.GetString("AI_API_URL"),
	}

	var provider ai.Provider
	switch *providerType{
	case "groq":
		request := ai.NewGroqRequest(aicfg)
		provider = ai.NewGroqClientWithRequest(aicfg, request)
		if len(systemPrompt.String()) > 0 {
			request.AddSystemMessage(systemPrompt.String())
		}
		request.AddMessage(prompt.String())
	default:
		panic("unknown provider type")
	}

	answer, err := provider.SendRequest(context.Background())
	if err != nil {
		panic(fmt.Errorf("error while asking ai: %w\n", err))
	}

	fmt.Println("PROMPT:")
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println(prompt.String())
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println()
	fmt.Println("ANSWER:")
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println(answer)
	fmt.Println(strings.Repeat("=", 20))
}
