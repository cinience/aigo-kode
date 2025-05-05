package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cinience/aigo-kode/internal/core"
	"github.com/sashabaranov/go-openai"
)

// OpenAIModel implements the AIModel interface for OpenAI
type OpenAIModel struct {
	client      *openai.Client
	modelName   string
	temperature float32
	maxTokens   int
}

// NewOpenAIModel creates a new OpenAI model
func NewOpenAIModel(apiKey, modelName, baseURL string) (*OpenAIModel, error) {
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}

	if modelName == "" {
		modelName = "gpt-3.5-turbo"
	}

	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	client := openai.NewClientWithConfig(config)

	return &OpenAIModel{
		client:      client,
		modelName:   modelName,
		temperature: 0.7,
		maxTokens:   4096,
	}, nil
}

// Name returns the model name
func (m *OpenAIModel) Name() string {
	return m.modelName
}

// Provider returns the model provider
func (m *OpenAIModel) Provider() string {
	return "OpenAI"
}

// Query sends a query to the model and returns a response
func (m *OpenAIModel) Query(ctx context.Context, messages []core.Message, tools []core.Tool) (*core.Response, error) {
	// Convert messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: fmt.Sprintf("%v", msg.Content),
		}
		if len(msg.ToolCalls) > 0 {
			openaiMessages[i].ToolCalls = make([]openai.ToolCall, len(msg.ToolCalls))
			for toolCallIdx, toolCall := range msg.ToolCalls {
				arguments, _ := json.Marshal(toolCall.Input)
				openaiMessages[i].ToolCalls[toolCallIdx] = openai.ToolCall{
					ID: toolCall.ID,
					Function: openai.FunctionCall{
						Name:      toolCall.ToolName,
						Arguments: string(arguments),
					},
				}
			}
		}
	}

	// Create request
	req := openai.ChatCompletionRequest{
		Model:       m.modelName,
		Messages:    openaiMessages,
		Temperature: m.temperature,
		MaxTokens:   m.maxTokens,
		Tools:       make([]openai.Tool, len(tools)),
	}
	for toolIdx, tool := range tools {
		req.Tools[toolIdx] = openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:       tool.Name(),
				Parameters: tool.Arguments(),
			},
		}
	}

	// Send request
	resp, err := m.client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	data, _ := json.Marshal(resp)
	log.Println("*************************")
	log.Println(string(data))
	log.Println("$$$$$$$$$$$$$$$$$$$$$$$")
	// Process response
	if len(resp.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	var content string
	var finishReason string

	// For simplicity, we're not handling tool calls in this version
	// since the API seems to have changed
	var toolCalls []core.ToolCall
	for _, choice := range resp.Choices {
		content = choice.Message.Content
		finishReason = string(choice.FinishReason)
		if choice.Message.ToolCalls != nil {
			for _, toolCall := range choice.Message.ToolCalls {
				input := make(map[string]interface{})
				_ = json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				toolCalls = append(toolCalls, core.ToolCall{
					ID:       toolCall.ID,
					ToolName: toolCall.Function.Name,
					Input:    input,
				})
			}
		}
	}

	return &core.Response{
		Content:   content,
		ToolCalls: toolCalls,
		Usage: core.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		FinishReason: finishReason,
	}, nil
}

// StreamQuery sends a query to the model and returns a stream of response chunks
func (m *OpenAIModel) StreamQuery(ctx context.Context, messages []core.Message, tools []core.Tool) (<-chan core.ResponseChunk, error) {
	// Convert messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: fmt.Sprintf("%v", msg.Content),
		}
	}

	// Create request
	req := openai.ChatCompletionRequest{
		Model:       m.modelName,
		Messages:    openaiMessages,
		Temperature: m.temperature,
		MaxTokens:   m.maxTokens,
		Stream:      true,
	}

	// Send request
	stream, err := m.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	// Create response channel
	responseCh := make(chan core.ResponseChunk)

	// Process stream in a goroutine
	go func() {
		defer close(responseCh)
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err != nil {
				if strings.Contains(err.Error(), "EOF") {
					// End of stream
					responseCh <- core.ResponseChunk{
						IsDone: true,
					}
					return
				}
				// Other error
				responseCh <- core.ResponseChunk{
					Error:  err,
					IsDone: true,
				}
				return
			}

			// Send chunk
			responseCh <- core.ResponseChunk{
				Content: response.Choices[0].Delta.Content,
				IsDone:  false,
			}
		}
	}()

	return responseCh, nil
}
