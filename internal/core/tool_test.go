package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToolInterface(t *testing.T) {
	// Create a mock tool for testing
	mockTool := &MockTool{
		name:        "MockTool",
		description: "A mock tool for testing",
		readOnly:    true,
	}

	// Test the basic interface methods
	assert.Equal(t, "MockTool", mockTool.Name())
	assert.Equal(t, "A mock tool for testing", mockTool.Description())
	assert.True(t, mockTool.IsReadOnly())
	assert.False(t, mockTool.RequiresPermission(nil))
}

// MockTool implements the Tool interface for testing
type MockTool struct {
	name        string
	description string
	readOnly    bool
	executeFunc func(map[string]interface{}) (interface{}, error)
}

func (t *MockTool) Name() string {
	return t.name
}

func (t *MockTool) Description() string {
	return t.description
}

func (t *MockTool) Execute(ctx interface{}, input map[string]interface{}) (interface{}, error) {
	if t.executeFunc != nil {
		return t.executeFunc(input)
	}
	return map[string]string{"result": "mock result"}, nil
}

func (t *MockTool) ValidateInput(input map[string]interface{}) error {
	return nil
}

func (t *MockTool) IsReadOnly() bool {
	return t.readOnly
}

func (t *MockTool) RequiresPermission(input map[string]interface{}) bool {
	return false
}
