// Package llm provides an LLM client abstraction for crossword generation.
package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// ErrMaxRetries is returned when max retries are exceeded.
var ErrMaxRetries = errors.New("max retries exceeded")

// Request represents an LLM request.
type Request struct {
	Prompt       string            `json:"prompt"`
	SystemPrompt string            `json:"system_prompt,omitempty"`
	MaxTokens    int               `json:"max_tokens,omitempty"`
	Temperature  float64           `json:"temperature,omitempty"`
	Schema       *jsonschema.Schema `json:"-"` // For output validation
	SchemaName   string            `json:"schema_name,omitempty"`
}

// Response represents an LLM response.
type Response struct {
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason,omitempty"`
	TokensUsed   int    `json:"tokens_used,omitempty"`
}

// Client is the interface for LLM providers.
type Client interface {
	// Complete sends a completion request to the LLM.
	Complete(ctx context.Context, req Request) (*Response, error)
}

// Config holds client configuration.
type Config struct {
	MaxRetries     int     // Max retry attempts for validation failures
	RepairPrompt   string  // Prompt template for repair attempts
	DefaultTemp    float64 // Default temperature
	DefaultTokens  int     // Default max tokens
	RedactSecrets  bool    // Whether to redact secrets in traces
}

// DefaultConfig returns default client configuration.
func DefaultConfig() Config {
	return Config{
		MaxRetries:    3,
		DefaultTemp:   0.7,
		DefaultTokens: 2048,
		RedactSecrets: true,
		RepairPrompt: `The previous response was invalid JSON or didn't match the required schema.
Error: %s
Previous response: %s

Please provide a corrected response that is valid JSON matching the schema.`,
	}
}

// ValidatingClient wraps a Client with JSON schema validation and retry logic.
type ValidatingClient struct {
	client Client
	config Config
	traces []Trace
}

// Trace records an LLM interaction for debugging.
type Trace struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Error    string   `json:"error,omitempty"`
	Attempt  int      `json:"attempt"`
}

// NewValidatingClient creates a new validating client wrapper.
func NewValidatingClient(client Client, config Config) *ValidatingClient {
	return &ValidatingClient{
		client: client,
		config: config,
	}
}

// CompleteWithValidation sends a request and validates the JSON response.
// It retries with repair prompts on validation failures.
func (c *ValidatingClient) CompleteWithValidation(ctx context.Context, req Request, target interface{}) error {
	var lastError error
	originalPrompt := req.Prompt

	for attempt := 1; attempt <= c.config.MaxRetries; attempt++ {
		resp, err := c.client.Complete(ctx, req)
		if err != nil {
			c.recordTrace(req, Response{}, err.Error(), attempt)
			return fmt.Errorf("LLM request failed: %w", err)
		}

		c.recordTrace(req, *resp, "", attempt)

		// Check for empty response
		if resp.Content == "" {
			lastError = fmt.Errorf("empty response from LLM (finish_reason: %s)", resp.FinishReason)
			continue
		}

		// Extract JSON from response (handle markdown code blocks)
		jsonContent := extractJSON(resp.Content)

		// Try to unmarshal
		if err := json.Unmarshal([]byte(jsonContent), target); err != nil {
			lastError = fmt.Errorf("JSON parse error: %w (response length: %d, content: %s)",
				err, len(resp.Content), truncate(resp.Content, 200))
			req.Prompt = fmt.Sprintf(c.config.RepairPrompt, lastError.Error(), truncate(resp.Content, 500))
			continue
		}

		// Validate against schema if provided
		if req.Schema != nil {
			var doc interface{}
			json.Unmarshal([]byte(jsonContent), &doc)
			if err := req.Schema.Validate(doc); err != nil {
				lastError = fmt.Errorf("schema validation error: %w", err)
				req.Prompt = fmt.Sprintf(c.config.RepairPrompt, lastError.Error(), truncate(resp.Content, 500))
				continue
			}
		}

		// Success
		return nil
	}

	// Restore original prompt for error context
	_ = originalPrompt
	return fmt.Errorf("%w: %v", ErrMaxRetries, lastError)
}

// Traces returns recorded traces (with secrets redacted if configured).
func (c *ValidatingClient) Traces() []Trace {
	if !c.config.RedactSecrets {
		return c.traces
	}

	redacted := make([]Trace, len(c.traces))
	for i, t := range c.traces {
		redacted[i] = Trace{
			Request: Request{
				Prompt:       redactSecrets(t.Request.Prompt),
				SystemPrompt: redactSecrets(t.Request.SystemPrompt),
				MaxTokens:    t.Request.MaxTokens,
				Temperature:  t.Request.Temperature,
				SchemaName:   t.Request.SchemaName,
			},
			Response: t.Response,
			Error:    t.Error,
			Attempt:  t.Attempt,
		}
	}
	return redacted
}

// ClearTraces clears recorded traces.
func (c *ValidatingClient) ClearTraces() {
	c.traces = nil
}

func (c *ValidatingClient) recordTrace(req Request, resp Response, errStr string, attempt int) {
	c.traces = append(c.traces, Trace{
		Request:  req,
		Response: resp,
		Error:    errStr,
		Attempt:  attempt,
	})
}

// extractJSON extracts JSON from a response that might be wrapped in markdown.
func extractJSON(content string) string {
	content = strings.TrimSpace(content)

	// Check for markdown code blocks
	if strings.HasPrefix(content, "```") {
		// Find the end of the first line (language specifier)
		firstNewline := strings.Index(content, "\n")
		if firstNewline > 0 {
			content = content[firstNewline+1:]
		}
		// Remove trailing ```
		if idx := strings.LastIndex(content, "```"); idx > 0 {
			content = content[:idx]
		}
	}

	return strings.TrimSpace(content)
}

// redactSecrets redacts potential secrets from text.
func redactSecrets(text string) string {
	// Redact API keys (common patterns)
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(api[_-]?key|apikey|secret|password|token)[=:]\s*["']?[a-zA-Z0-9_-]{10,}["']?`),
		regexp.MustCompile(`sk-ant-[a-zA-Z0-9-]{10,}`),      // Anthropic keys (match first)
		regexp.MustCompile(`sk-[a-zA-Z0-9-]{10,}`),          // OpenAI keys
		regexp.MustCompile(`Bearer\s+[a-zA-Z0-9._-]{10,}`),  // Bearer tokens
	}

	result := text
	for _, p := range patterns {
		result = p.ReplaceAllString(result, "[REDACTED]")
	}
	return result
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
