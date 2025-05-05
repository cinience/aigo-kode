package core

import (
	"context"
)

// Message represents a message in a conversation
type Message struct {
	// Role can be "user", "assistant", or "system"
	Role string
	// Content can be a string or structured content
	Content interface{}

	// ToolCalls contains any tool calls in the response
	ToolCalls []ToolCall
}

// ToolCall represents a request from the AI to use a tool
type ToolCall struct {
	// ToolName is the name of the tool to call
	ToolName string
	// Input contains the parameters for the tool
	Input map[string]interface{}
	// ID is a unique identifier for this tool call
	ID string
}

// Usage represents API usage statistics
type Usage struct {
	// PromptTokens is the number of tokens in the prompt
	PromptTokens int
	// CompletionTokens is the number of tokens in the completion
	CompletionTokens int
	// TotalTokens is the total number of tokens
	TotalTokens int
	// Cost is the estimated cost in USD
	Cost float64
}

// ResponseChunk represents a chunk of a streaming response
type ResponseChunk struct {
	// Content is the text content of this chunk
	Content string
	// ToolCalls contains any tool calls in this chunk
	ToolCalls []ToolCall
	// IsDone indicates if this is the final chunk
	IsDone bool
	// Error contains any error that occurred
	Error error
}

// Response represents a complete model response
type Response struct {
	// Content is the text content of the response
	Content string
	// ToolCalls contains any tool calls in the response
	ToolCalls []ToolCall
	// Usage contains token usage statistics
	Usage Usage
	// FinishReason indicates why the model stopped generating
	FinishReason string
}

// AIModel defines the interface for AI model providers
type AIModel interface {
	// Query sends a query to the model and returns a response
	Query(ctx context.Context, messages []Message, tools []Tool) (*Response, error)

	// StreamQuery sends a query to the model and returns a stream of response chunks
	StreamQuery(ctx context.Context, messages []Message, tools []Tool) (<-chan ResponseChunk, error)

	// Name returns the model name
	Name() string

	// Provider returns the model provider
	Provider() string
}
