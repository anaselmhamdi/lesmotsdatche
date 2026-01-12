package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIConfig holds OpenAI-specific configuration.
type OpenAIConfig struct {
	APIKey       string
	Model        string
	BaseURL      string
	Timeout      time.Duration
	Organization string
}

// DefaultOpenAIConfig returns default OpenAI configuration.
func DefaultOpenAIConfig() OpenAIConfig {
	return OpenAIConfig{
		Model:   "gpt-4o",
		BaseURL: "https://api.openai.com/v1",
		Timeout: 60 * time.Second,
	}
}

// OpenAIClient implements the Client interface for OpenAI's API.
type OpenAIClient struct {
	config     OpenAIConfig
	httpClient *http.Client
}

// NewOpenAIClient creates a new OpenAI client.
func NewOpenAIClient(config OpenAIConfig) *OpenAIClient {
	if config.BaseURL == "" {
		config.BaseURL = DefaultOpenAIConfig().BaseURL
	}
	if config.Model == "" {
		config.Model = DefaultOpenAIConfig().Model
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultOpenAIConfig().Timeout
	}

	return &OpenAIClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// openAIRequest is the request structure for OpenAI's chat completions API.
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse is the response structure from OpenAI's chat completions API.
type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      openAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *openAIError `json:"error,omitempty"`
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// Complete sends a completion request to OpenAI.
func (c *OpenAIClient) Complete(ctx context.Context, req Request) (*Response, error) {
	messages := []openAIMessage{}

	if req.SystemPrompt != "" {
		messages = append(messages, openAIMessage{
			Role:    "system",
			Content: req.SystemPrompt,
		})
	}

	messages = append(messages, openAIMessage{
		Role:    "user",
		Content: req.Prompt,
	})

	openaiReq := openAIRequest{
		Model:       c.config.Model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	if openaiReq.MaxTokens == 0 {
		openaiReq.MaxTokens = 2048
	}
	if openaiReq.Temperature == 0 {
		openaiReq.Temperature = 0.7
	}

	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	if c.config.Organization != "" {
		httpReq.Header.Set("OpenAI-Organization", c.config.Organization)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var openaiResp openAIResponse
	if err := json.Unmarshal(respBody, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if openaiResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s (type: %s, code: %s)",
			openaiResp.Error.Message, openaiResp.Error.Type, openaiResp.Error.Code)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &Response{
		Content:      openaiResp.Choices[0].Message.Content,
		FinishReason: openaiResp.Choices[0].FinishReason,
		TokensUsed:   openaiResp.Usage.TotalTokens,
	}, nil
}

// Provider returns the provider name.
func (c *OpenAIClient) Provider() string {
	return "openai"
}

// Model returns the model name being used.
func (c *OpenAIClient) Model() string {
	return c.config.Model
}
