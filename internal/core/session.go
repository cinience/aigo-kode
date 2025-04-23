package core

import (
	"context"
	"errors"
)

// Session represents an interactive session with an AI model
type Session struct {
	// Messages is the conversation history
	Messages []Message
	// Model is the AI model being used
	Model AIModel
	// Tools are the available tools
	Tools []Tool
	// Config is the session configuration
	Config *SessionConfig
}

// SessionConfig contains configuration for a session
type SessionConfig struct {
	// ProjectPath is the path to the current project
	ProjectPath string
	// SystemPrompt is the system prompt to use
	SystemPrompt string
	// MaxTokens is the maximum number of tokens to generate
	MaxTokens int
	// Temperature controls randomness (0.0-2.0)
	Temperature float64
}

// NewSession creates a new session with the given model and tools
func NewSession(model AIModel, tools []Tool, config *SessionConfig) *Session {
	if config == nil {
		config = &SessionConfig{
			MaxTokens:    4096,
			Temperature:  0.7,
			SystemPrompt: "You are a helpful AI assistant that can use tools to help with coding tasks.",
		}
	}

	return &Session{
		Messages: []Message{
			{
				Role:    "system",
				Content: config.SystemPrompt,
			},
		},
		Model:  model,
		Tools:  tools,
		Config: config,
	}
}

// AddUserMessage adds a user message to the conversation
func (s *Session) AddUserMessage(content string) {
	s.Messages = append(s.Messages, Message{
		Role:    "user",
		Content: content,
	})
}

// AddAssistantMessage adds an assistant message to the conversation
func (s *Session) AddAssistantMessage(content string) {
	s.Messages = append(s.Messages, Message{
		Role:    "assistant",
		Content: content,
	})
}

// AddToolResult adds a tool result to the conversation
func (s *Session) AddToolResult(toolName string, id string, result interface{}) {
	// In a real implementation, this would format the result properly
	// based on the specific tool and result type
	s.Messages = append(s.Messages, Message{
		Role:    "tool",
		Content: result,
	})
}

// Query sends the current conversation to the model and returns a response
func (s *Session) Query(ctx context.Context) (*Response, error) {
	if s.Model == nil {
		return nil, errors.New("no model configured")
	}

	return s.Model.Query(ctx, s.Messages, s.Tools)
}

// StreamQuery sends the current conversation to the model and returns a stream of response chunks
func (s *Session) StreamQuery(ctx context.Context) (<-chan ResponseChunk, error) {
	if s.Model == nil {
		return nil, errors.New("no model configured")
	}

	return s.Model.StreamQuery(ctx, s.Messages, s.Tools)
}

// ExecuteTool executes a tool and adds the result to the conversation
func (s *Session) ExecuteTool(ctx context.Context, toolCall ToolCall) (*ToolUseResult, error) {
	// Find the tool
	var tool Tool
	for _, t := range s.Tools {
		if t.Name() == toolCall.ToolName {
			tool = t
			break
		}
	}

	if tool == nil {
		return nil, errors.New("tool not found: " + toolCall.ToolName)
	}

	// Check if tool requires permission
	if tool.RequiresPermission(toolCall.Input) {
		// In a real implementation, this would prompt the user for permission
		// For now, we'll just allow it
	}

	// Validate input
	if err := tool.ValidateInput(toolCall.Input); err != nil {
		return &ToolUseResult{
			ToolName: toolCall.ToolName,
			Input:    toolCall.Input,
			Error:    err,
		}, nil
	}

	// Execute tool
	output, err := tool.Execute(ctx, toolCall.Input)
	result := &ToolUseResult{
		ToolName: toolCall.ToolName,
		Input:    toolCall.Input,
		Output:   output,
		Error:    err,
	}

	// Add result to conversation
	s.AddToolResult(toolCall.ToolName, toolCall.ID, output)

	return result, nil
}
