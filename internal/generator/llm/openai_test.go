package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOpenAIClient_Complete(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected /chat/completions, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("expected Bearer test-key, got %s", r.Header.Get("Authorization"))
		}

		// Decode request body
		var req openAIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Model != "gpt-4o" {
			t.Errorf("expected model gpt-4o, got %s", req.Model)
		}
		if len(req.Messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(req.Messages))
		}

		// Send mock response
		resp := openAIResponse{
			ID:      "chatcmpl-test",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "gpt-4o",
			Choices: []struct {
				Index        int           `json:"index"`
				Message      openAIMessage `json:"message"`
				FinishReason string        `json:"finish_reason"`
			}{
				{
					Index:        0,
					Message:      openAIMessage{Role: "assistant", Content: `{"result": "test"}`},
					FinishReason: "stop",
				},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewOpenAIClient(OpenAIConfig{
		APIKey:  "test-key",
		Model:   "gpt-4o",
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	resp, err := client.Complete(context.Background(), Request{
		SystemPrompt: "You are a test assistant",
		Prompt:       "Test prompt",
		MaxTokens:    100,
		Temperature:  0.7,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Content != `{"result": "test"}` {
		t.Errorf("unexpected content: %s", resp.Content)
	}
	if resp.FinishReason != "stop" {
		t.Errorf("expected finish_reason 'stop', got %s", resp.FinishReason)
	}
	if resp.TokensUsed != 15 {
		t.Errorf("expected 15 tokens, got %d", resp.TokensUsed)
	}
}

func TestOpenAIClient_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := openAIResponse{
			Error: &openAIError{
				Message: "Invalid API key",
				Type:    "invalid_request_error",
				Code:    "invalid_api_key",
			},
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewOpenAIClient(OpenAIConfig{
		APIKey:  "invalid-key",
		BaseURL: server.URL,
	})

	_, err := client.Complete(context.Background(), Request{
		Prompt: "Test",
	})

	if err == nil {
		t.Error("expected error for invalid API key")
	}
}

func TestOpenAIClient_NoChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := openAIResponse{
			Choices: []struct {
				Index        int           `json:"index"`
				Message      openAIMessage `json:"message"`
				FinishReason string        `json:"finish_reason"`
			}{},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewOpenAIClient(OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := client.Complete(context.Background(), Request{
		Prompt: "Test",
	})

	if err == nil {
		t.Error("expected error for no choices")
	}
}

func TestOpenAIClient_DefaultConfig(t *testing.T) {
	config := DefaultOpenAIConfig()

	if config.Model != "gpt-4o" {
		t.Errorf("expected default model gpt-4o, got %s", config.Model)
	}
	if config.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("unexpected base URL: %s", config.BaseURL)
	}
	if config.Timeout != 60*time.Second {
		t.Errorf("expected 60s timeout, got %v", config.Timeout)
	}
}

func TestOpenAIClient_ProviderInfo(t *testing.T) {
	client := NewOpenAIClient(OpenAIConfig{
		Model: "gpt-4o-mini",
	})

	if client.Provider() != "openai" {
		t.Errorf("expected provider 'openai', got %s", client.Provider())
	}
	if client.Model() != "gpt-4o-mini" {
		t.Errorf("expected model 'gpt-4o-mini', got %s", client.Model())
	}
}

func TestOpenAIClient_Organization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		org := r.Header.Get("OpenAI-Organization")
		if org != "test-org" {
			t.Errorf("expected organization 'test-org', got %s", org)
		}

		resp := openAIResponse{
			Choices: []struct {
				Index        int           `json:"index"`
				Message      openAIMessage `json:"message"`
				FinishReason string        `json:"finish_reason"`
			}{
				{Message: openAIMessage{Content: "ok"}, FinishReason: "stop"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewOpenAIClient(OpenAIConfig{
		APIKey:       "test-key",
		BaseURL:      server.URL,
		Organization: "test-org",
	})

	_, err := client.Complete(context.Background(), Request{Prompt: "Test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
