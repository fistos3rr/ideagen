package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GroqClient struct {
	cfg  Config
	http *http.Client
	Req  *GroqRequest
}

func NewGroqClient(cfg Config) *GroqClient {
	return &GroqClient{
		cfg:  cfg,
		http: &http.Client{Timeout: 30 * time.Second},
		Req:  NewGroqRequest(cfg),
	}
}

func NewGroqClientWithRequest(cfg Config, req *GroqRequest) *GroqClient {
	return &GroqClient{
		cfg:  cfg,
		http: &http.Client{Timeout: 30 * time.Second},
		Req:  req,
	}
}

type GroqRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	TopP        float64       `json:"top_p"`
}

func NewGroqRequest(cfg Config) *GroqRequest {
	return &GroqRequest{
		Model:       cfg.Model,
		Messages:    make([]groqMessage, 0),
		Temperature: 1.0,
		TopP:        1.0,
	}
}

func (req *GroqRequest) Clear(message string) {
	req.Messages = make([]groqMessage, 0)
}

func (req *GroqRequest) AddMessage(message string) {
	req.Messages = append(req.Messages, groqMessage{Role: "user", Content: message})
}

func (req *GroqRequest) AddSystemMessage(message string) {
	req.Messages = append(req.Messages, groqMessage{Role: "system", Content: message})
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *GroqClient) SendMessage(ctx context.Context, message string) (string, error) {
	c.Req.Messages = []groqMessage{
		{Role: "user", Content: message},
	}

	jsonData, err := json.Marshal(c.Req)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.cfg.APIURL, bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq API returned status %d", resp.StatusCode)
	}

	var result groqResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return result.Choices[0].Message.Content, nil
}

func (c *GroqClient) SendRequest(ctx context.Context) (string, error) {
	jsonData, err := json.Marshal(c.Req)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.cfg.APIURL, bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq API returned status %d", resp.StatusCode)
	}

	var result groqResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return result.Choices[0].Message.Content, nil
}
