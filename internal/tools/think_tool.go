package tools

import (
	"context"
	"errors"

	"github.com/cinience/aigo-kode/internal/core"
)

// ThinkTool implements the Tool interface for AI thinking/reasoning
type ThinkTool struct{}

// Name returns the tool name
func (t *ThinkTool) Name() string {
	return "Think"
}

// Description returns the tool description
func (t *ThinkTool) Description() string {
	return "Allows the AI to think through a problem step by step"
}

// ThinkToolOutput defines the output structure for ThinkTool
type ThinkToolOutput struct {
	Reasoning string `json:"reasoning"`
}

// Execute executes the thinking process
func (t *ThinkTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	// Extract prompt
	prompt, ok := input["prompt"].(string)
	if !ok || prompt == "" {
		return nil, errors.New("prompt is required and must be a string")
	}

	// This tool doesn't actually do anything except return the prompt
	// It's meant to be used by the AI to think through problems
	return &ThinkToolOutput{
		Reasoning: prompt,
	}, nil
}

// ValidateInput validates the input parameters
func (t *ThinkTool) ValidateInput(input map[string]interface{}) error {
	// Check if prompt exists and is a string
	promptVal, ok := input["prompt"]
	if !ok {
		return errors.New("prompt is required")
	}

	prompt, ok := promptVal.(string)
	if !ok {
		return errors.New("prompt must be a string")
	}

	if prompt == "" {
		return errors.New("prompt cannot be empty")
	}

	return nil
}

func (t *ThinkTool) Arguments() string {
	return `{
		"prompt": {
			"type": "string",
			"description": "The prompt to think through"
		}
	}`
}

// IsReadOnly returns whether the tool is read-only
func (t *ThinkTool) IsReadOnly() bool {
	return true
}

// RequiresPermission checks if permission is needed
func (t *ThinkTool) RequiresPermission(input map[string]interface{}) bool {
	return false
}

// NewThinkTool creates a new ThinkTool
func NewThinkTool() core.Tool {
	return &ThinkTool{}
}
