package llm

import (
	"context"
	"errors"
)

// MockClient is a mock LLM client for testing.
type MockClient struct {
	Responses []string // Responses to return in order
	Errors    []error  // Errors to return in order
	Calls     []Request // Recorded calls
	callIndex int
}

// NewMockClient creates a new mock client.
func NewMockClient(responses ...string) *MockClient {
	return &MockClient{Responses: responses}
}

// WithErrors sets errors to return.
func (m *MockClient) WithErrors(errs ...error) *MockClient {
	m.Errors = errs
	return m
}

// Complete returns the next mock response.
func (m *MockClient) Complete(ctx context.Context, req Request) (*Response, error) {
	m.Calls = append(m.Calls, req)

	// Check for error
	if m.callIndex < len(m.Errors) && m.Errors[m.callIndex] != nil {
		err := m.Errors[m.callIndex]
		m.callIndex++
		return nil, err
	}

	// Return response
	if m.callIndex >= len(m.Responses) {
		m.callIndex++
		return nil, errors.New("no more mock responses")
	}

	resp := &Response{
		Content:      m.Responses[m.callIndex],
		FinishReason: "stop",
		TokensUsed:   100,
	}
	m.callIndex++
	return resp, nil
}

// Reset resets the mock client state.
func (m *MockClient) Reset() {
	m.callIndex = 0
	m.Calls = nil
}

// CallCount returns the number of calls made.
func (m *MockClient) CallCount() int {
	return len(m.Calls)
}
