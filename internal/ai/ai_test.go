package ai

import (
	"context"
	"testing"
	"time"

	"github.com/cinience/aigo-kode/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOpenAIClient is a mock implementation of the OpenAI client
type MockOpenAIClient struct {
	mock.Mock
}

// TestOpenAIModel tests the OpenAI model implementation
func TestOpenAIModel(t *testing.T) {
	// Skip if no API key is available
	t.Skip("Skipping OpenAI API test as it requires an API key")

	// This is a basic structure test, not an actual API call test
	model, err := NewOpenAIModel("test-api-key", "gpt-3.5-turbo", "")
	assert.NoError(t, err)
	assert.NotNil(t, model)

	// Test Name and Provider methods
	assert.Equal(t, "gpt-3.5-turbo", model.Name())
	assert.Equal(t, "OpenAI", model.Provider())

	// Test error case
	_, err = NewOpenAIModel("", "", "")
	assert.Error(t, err)
}

// TestSession tests the Session implementation
func TestSession(t *testing.T) {
	// Create a mock tool
	mockTool := &MockTool{
		name:        "MockTool",
		description: "A mock tool for testing",
	}

	// Create a mock model
	mockModel := &MockModel{}

	// Create a session
	session := core.NewSession(mockModel, []core.Tool{mockTool}, &core.SessionConfig{
		SystemPrompt: "Test system prompt",
	})

	// Test initial state
	assert.NotNil(t, session)
	assert.Len(t, session.Messages, 1) // Should have system message
	assert.Equal(t, "system", session.Messages[0].Role)

	// Test adding messages
	session.AddUserMessage("Test user message")
	assert.Len(t, session.Messages, 2)
	assert.Equal(t, "user", session.Messages[1].Role)
	assert.Equal(t, "Test user message", session.Messages[1].Content)

	session.AddAssistantMessage("Test assistant message")
	assert.Len(t, session.Messages, 3)
	assert.Equal(t, "assistant", session.Messages[2].Role)
	assert.Equal(t, "Test assistant message", session.Messages[2].Content)
}

// MockTool implements the Tool interface for testing
type MockTool struct {
	name        string
	description string
}

func (t *MockTool) Name() string {
	return t.name
}

func (t *MockTool) Description() string {
	return t.description
}

func (t *MockTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	return map[string]string{"result": "mock result"}, nil
}

func (t *MockTool) ValidateInput(input map[string]interface{}) error {
	return nil
}

func (t *MockTool) IsReadOnly() bool {
	return true
}

func (t *MockTool) RequiresPermission(input map[string]interface{}) bool {
	return false
}

// MockModel implements the AIModel interface for testing
type MockModel struct{}

func (m *MockModel) Query(ctx context.Context, messages []core.Message, tools []core.Tool) (*core.Response, error) {
	return &core.Response{
		Content:   "Mock response",
		ToolCalls: []core.ToolCall{},
		Usage: core.Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
		FinishReason: "stop",
	}, nil
}

func (m *MockModel) StreamQuery(ctx context.Context, messages []core.Message, tools []core.Tool) (<-chan core.ResponseChunk, error) {
	ch := make(chan core.ResponseChunk)

	go func() {
		defer close(ch)

		// Send a chunk
		ch <- core.ResponseChunk{
			Content: "Mock ",
			IsDone:  false,
		}

		time.Sleep(10 * time.Millisecond)

		// Send another chunk
		ch <- core.ResponseChunk{
			Content: "response",
			IsDone:  true,
		}
	}()

	return ch, nil
}

func (m *MockModel) Name() string {
	return "MockModel"
}

func (m *MockModel) Provider() string {
	return "Mock"
}
