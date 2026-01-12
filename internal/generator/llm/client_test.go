package llm

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestValidatingClient_Success(t *testing.T) {
	mock := NewMockClient(`{"name": "test", "value": 42}`)
	client := NewValidatingClient(mock, DefaultConfig())

	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	err := client.CompleteWithValidation(context.Background(), Request{
		Prompt: "Generate JSON",
	}, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "test" {
		t.Errorf("expected name 'test', got %s", result.Name)
	}
	if result.Value != 42 {
		t.Errorf("expected value 42, got %d", result.Value)
	}

	if mock.CallCount() != 1 {
		t.Errorf("expected 1 call, got %d", mock.CallCount())
	}
}

func TestValidatingClient_MarkdownCodeBlock(t *testing.T) {
	mock := NewMockClient("```json\n{\"name\": \"test\"}\n```")
	client := NewValidatingClient(mock, DefaultConfig())

	var result struct {
		Name string `json:"name"`
	}

	err := client.CompleteWithValidation(context.Background(), Request{
		Prompt: "Generate JSON",
	}, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "test" {
		t.Errorf("expected name 'test', got %s", result.Name)
	}
}

func TestValidatingClient_RetryOnInvalidJSON(t *testing.T) {
	// First response is invalid, second is valid
	mock := NewMockClient(
		"not valid json",
		`{"name": "fixed"}`,
	)
	client := NewValidatingClient(mock, DefaultConfig())

	var result struct {
		Name string `json:"name"`
	}

	err := client.CompleteWithValidation(context.Background(), Request{
		Prompt: "Generate JSON",
	}, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "fixed" {
		t.Errorf("expected name 'fixed', got %s", result.Name)
	}

	if mock.CallCount() != 2 {
		t.Errorf("expected 2 calls (retry), got %d", mock.CallCount())
	}

	// Second call should be a repair prompt
	if !strings.Contains(mock.Calls[1].Prompt, "invalid JSON") {
		t.Error("expected repair prompt to mention invalid JSON")
	}
}

func TestValidatingClient_MaxRetriesExceeded(t *testing.T) {
	// All responses are invalid
	mock := NewMockClient(
		"invalid 1",
		"invalid 2",
		"invalid 3",
	)
	config := DefaultConfig()
	config.MaxRetries = 3
	client := NewValidatingClient(mock, config)

	var result struct {
		Name string `json:"name"`
	}

	err := client.CompleteWithValidation(context.Background(), Request{
		Prompt: "Generate JSON",
	}, &result)

	if !errors.Is(err, ErrMaxRetries) {
		t.Errorf("expected ErrMaxRetries, got: %v", err)
	}

	if mock.CallCount() != 3 {
		t.Errorf("expected 3 calls, got %d", mock.CallCount())
	}
}

func TestValidatingClient_LLMError(t *testing.T) {
	mock := NewMockClient().WithErrors(errors.New("API error"))
	client := NewValidatingClient(mock, DefaultConfig())

	var result struct{}
	err := client.CompleteWithValidation(context.Background(), Request{
		Prompt: "Generate JSON",
	}, &result)

	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "API error") {
		t.Errorf("expected API error in message, got: %v", err)
	}
}

func TestValidatingClient_Traces(t *testing.T) {
	mock := NewMockClient(`{"name": "test"}`)
	client := NewValidatingClient(mock, DefaultConfig())

	var result struct {
		Name string `json:"name"`
	}

	client.CompleteWithValidation(context.Background(), Request{
		Prompt: "Test prompt",
	}, &result)

	traces := client.Traces()
	if len(traces) != 1 {
		t.Fatalf("expected 1 trace, got %d", len(traces))
	}

	if traces[0].Attempt != 1 {
		t.Errorf("expected attempt 1, got %d", traces[0].Attempt)
	}
}

func TestValidatingClient_TracesRedaction(t *testing.T) {
	mock := NewMockClient(`{"name": "test"}`)
	config := DefaultConfig()
	config.RedactSecrets = true
	client := NewValidatingClient(mock, config)

	var result struct {
		Name string `json:"name"`
	}

	client.CompleteWithValidation(context.Background(), Request{
		Prompt: "Use API key: sk-ant-abc123xyz789def456",
	}, &result)

	traces := client.Traces()
	if strings.Contains(traces[0].Request.Prompt, "sk-ant-") {
		t.Error("expected API key to be redacted")
	}
	if !strings.Contains(traces[0].Request.Prompt, "[REDACTED]") {
		t.Error("expected [REDACTED] placeholder")
	}
}

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain JSON",
			input:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
		{
			name:     "markdown code block",
			input:    "```json\n{\"key\": \"value\"}\n```",
			expected: `{"key": "value"}`,
		},
		{
			name:     "markdown without language",
			input:    "```\n{\"key\": \"value\"}\n```",
			expected: `{"key": "value"}`,
		},
		{
			name:     "with whitespace",
			input:    "  \n{\"key\": \"value\"}\n  ",
			expected: `{"key": "value"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractJSON(tc.input)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestRedactSecrets(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string // Should not contain
	}{
		{
			name:     "OpenAI key",
			input:    "Use key: sk-abc123def456xyz789012345",
			contains: "sk-abc",
		},
		{
			name:     "Anthropic key",
			input:    "API: sk-ant-abc123def456xyz789012345",
			contains: "sk-ant-",
		},
		{
			name:     "Bearer token",
			input:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			contains: "eyJ",
		},
		{
			name:     "api_key param",
			input:    "api_key=secretkey12345678901234567890",
			contains: "secretkey",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := redactSecrets(tc.input)
			if strings.Contains(result, tc.contains) {
				t.Errorf("secret not redacted: %s", result)
			}
			if !strings.Contains(result, "[REDACTED]") {
				t.Errorf("expected [REDACTED] placeholder in: %s", result)
			}
		})
	}
}

func TestMockClient(t *testing.T) {
	mock := NewMockClient("response 1", "response 2")

	resp1, err := mock.Complete(context.Background(), Request{Prompt: "prompt 1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp1.Content != "response 1" {
		t.Errorf("expected 'response 1', got %s", resp1.Content)
	}

	resp2, err := mock.Complete(context.Background(), Request{Prompt: "prompt 2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp2.Content != "response 2" {
		t.Errorf("expected 'response 2', got %s", resp2.Content)
	}

	if mock.CallCount() != 2 {
		t.Errorf("expected 2 calls, got %d", mock.CallCount())
	}
}
