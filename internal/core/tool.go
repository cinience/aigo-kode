package core

import (
	"context"
	"errors"
)

// Tool defines the interface that all tools must implement
type Tool interface {
	// Name returns the tool name
	Name() string

	// Description returns the tool description
	Description() string

	Arguments() string

	// Execute executes the tool and returns the result
	Execute(ctx context.Context, input map[string]interface{}) (interface{}, error)

	// ValidateInput validates the input parameters
	ValidateInput(input map[string]interface{}) error

	// IsReadOnly returns whether the tool is read-only
	IsReadOnly() bool

	// RequiresPermission checks if permission is needed
	RequiresPermission(input map[string]interface{}) bool
}

// ToolUseResult represents the result of a tool execution
type ToolUseResult struct {
	ToolName string
	Input    map[string]interface{}
	Output   interface{}
	Error    error
}

// ErrPermissionDenied is returned when permission to use a tool is denied
var ErrPermissionDenied = errors.New("permission denied")

// ErrInvalidInput is returned when tool input validation fails
var ErrInvalidInput = errors.New("invalid input")

// ErrToolExecutionFailed is returned when tool execution fails
var ErrToolExecutionFailed = errors.New("tool execution failed")
