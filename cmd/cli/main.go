package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	"github.com/cinience/aigo-kode/internal/ai"
	"github.com/cinience/aigo-kode/internal/config"
	"github.com/cinience/aigo-kode/internal/core"
	"github.com/cinience/aigo-kode/internal/tools"
)

// Model represents the application state
type Model struct {
	session      *core.Session
	input        string
	messages     []string
	currentTool  string
	toolRegistry *tools.ToolRegistry
	config       *config.FileConfig
	err          error
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and user input
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.input == "" {
				return m, nil
			}

			// Add user message
			m.session.AddUserMessage(m.input)
			m.messages = append(m.messages, fmt.Sprintf("User: %s", m.input))

			// Clear input
			userInput := m.input
			m.input = ""

			// Return command to process the message
			return m, func() tea.Msg {
				return userMessageMsg(userInput)
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil
		default:
			m.input += msg.String()
			return m, nil
		}
	case userMessageMsg:
		// Process user message
		ctx := context.Background()
		resp, err := m.session.Query(ctx)
		if err != nil {
			m.err = err
			return m, nil
		}

		// Handle response
		if len(resp.ToolCalls) > 0 {
			// Tool calls
			for _, toolCall := range resp.ToolCalls {
				m.currentTool = toolCall.ToolName
				m.messages = append(m.messages, fmt.Sprintf("AI wants to use tool: %s", toolCall.ToolName))

				// Execute tool
				result, err := m.session.ExecuteTool(ctx, toolCall)
				if err != nil {
					m.err = err
					return m, nil
				}

				if result.Error != nil {
					m.messages = append(m.messages, fmt.Sprintf("Tool error: %v", result.Error))
				} else {
					m.messages = append(m.messages, fmt.Sprintf("Tool result: %v", result.Output))
				}
			}

			// Get final response after tool use
			resp, err = m.session.Query(ctx)
			if err != nil {
				m.err = err
				return m, nil
			}
		}

		// Add assistant message
		m.session.AddAssistantMessage(resp.Content)
		m.messages = append(m.messages, fmt.Sprintf("AI: %s", resp.Content))
		m.currentTool = ""

		return m, nil
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	// Simple terminal UI
	view := "Go Anon Kode - Terminal AI Assistant\n\n"

	// Show messages
	for _, msg := range m.messages {
		view += msg + "\n\n"
	}

	// Show current tool if any
	if m.currentTool != "" {
		view += fmt.Sprintf("Using tool: %s\n", m.currentTool)
	}

	// Show error if any
	if m.err != nil {
		view += fmt.Sprintf("Error: %v\n", m.err)
	}

	// Show input prompt
	view += "\n> " + m.input

	return view
}

// userMessageMsg represents a user message
type userMessageMsg string

func main() {
	// Set up config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".go-anon-kode")
	cfg, err := config.NewFileConfig(configDir)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	// Get global config
	globalConfig, err := cfg.GetGlobalConfig()
	if err != nil {
		log.Fatalf("Failed to get global config: %v", err)
	}

	// Set up OpenAI model
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// Try to get from config
		apiKey = globalConfig.APIKeys["openai"]
		if apiKey == "" {
			log.Fatalf("OpenAI API key not found. Set OPENAI_API_KEY environment variable or configure in settings.")
		}
	}

	model, err := ai.NewOpenAIModel(apiKey, globalConfig.DefaultModel)
	if err != nil {
		log.Fatalf("Failed to create OpenAI model: %v", err)
	}

	// Set up tool registry
	registry := tools.DefaultToolRegistry()

	// Create session
	session := core.NewSession(model, registry.GetAllTools(), &core.SessionConfig{
		ProjectPath:  ".",
		SystemPrompt: "You are a helpful AI coding assistant. You can help with coding tasks, answer questions, and use tools to interact with the file system.",
		MaxTokens:    4096,
		Temperature:  0.7,
	})

	// Create and run the Bubble Tea application
	p := tea.NewProgram(Model{
		session:      session,
		toolRegistry: registry,
		config:       cfg,
		messages:     []string{"Welcome to Go Anon Kode! Type your question or request."},
	})

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
