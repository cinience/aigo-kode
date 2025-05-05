package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cinience/aigo-kode/internal/ai"
	"github.com/cinience/aigo-kode/internal/config"
	"github.com/cinience/aigo-kode/internal/core"
	"github.com/cinience/aigo-kode/internal/tools"
	"github.com/gin-gonic/gin"
)

// Server represents the web server
type Server struct {
	router       *gin.Engine
	toolRegistry *tools.ToolRegistry
	config       *config.FileConfig
}

// SessionStore manages active sessions
type SessionStore struct {
	sessions map[string]*core.Session
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*core.Session),
	}
}

// GetSession retrieves a session by ID, creating it if it doesn't exist
func (s *SessionStore) GetSession(id string, model core.AIModel, tools []core.Tool) *core.Session {
	if session, ok := s.sessions[id]; ok {
		return session
	}

	// Create new session
	session := core.NewSession(model, tools, &core.SessionConfig{
		ProjectPath:  ".",
		SystemPrompt: "You are a helpful AI coding assistant. You can help with coding tasks, answer questions, and use tools to interact with the file system.",
		MaxTokens:    4096,
		Temperature:  0.7,
	})

	s.sessions[id] = session
	return session
}

// NewServer creates a new web server
func NewServer() (*Server, error) {
	// Set up config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".go-anon-kode")
	cfg, err := config.NewFileConfig(configDir)
	if err != nil {
		return nil, err
	}

	// Set up tool registry
	registry := tools.DefaultToolRegistry()

	// Create router
	router := gin.Default()

	// Create server
	server := &Server{
		router:       router,
		toolRegistry: registry,
		config:       cfg,
	}

	// Set up routes
	server.setupRoutes()

	return server, nil
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Serve static files for the web UI
	s.router.Static("/assets", "./web/dist/assets")
	s.router.StaticFile("/", "./web/dist/index.html")
	s.router.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

	// API routes
	api := s.router.Group("/api")
	{
		api.POST("/chat", s.handleChat)
		api.GET("/chat/history", s.handleGetChatHistory)
		api.POST("/tools/:toolName", s.handleExecuteTool)
		api.GET("/files", s.handleListFiles)
		api.GET("/files/:path", s.handleGetFile)
		api.PUT("/files/:path", s.handleUpdateFile)
		api.GET("/config", s.handleGetConfig)
		api.PUT("/config", s.handleUpdateConfig)
	}
}

// Message represents a chat message
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ChatRequest represents a chat API request
type ChatRequest struct {
	SessionID string `json:"sessionId"`
	Message   string `json:"message"`
}

// ChatResponse represents a chat API response
type ChatResponse struct {
	Response  string          `json:"response"`
	ToolCalls []core.ToolCall `json:"toolCalls,omitempty"`
}

// handleChat handles chat API requests
func (s *Server) handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get global config
	globalConfig, err := s.config.GetGlobalConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get config"})
		return
	}

	// Set up OpenAI model
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// Try to get from config
		apiKey = globalConfig.APIKeys["openai"]
		if apiKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI API key not found"})
			return
		}
	}

	model, err := ai.NewOpenAIModel(apiKey, globalConfig.DefaultModel, globalConfig.BaseURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create model"})
		return
	}

	// Get or create session
	sessionStore := NewSessionStore()
	session := sessionStore.GetSession(req.SessionID, model, s.toolRegistry.GetAllTools())

	// Add user message
	session.AddUserMessage(req.Message)

	// Query model
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := session.Query(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle tool calls if any
	if len(resp.ToolCalls) > 0 {
		for _, toolCall := range resp.ToolCalls {
			// Execute tool
			result, err := session.ExecuteTool(ctx, toolCall)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if result.Error != nil {
				// Tool execution failed, but we continue
				log.Printf("Tool execution failed: %v", result.Error)
			}
		}

		// Get final response after tool use
		resp, err = session.Query(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Add assistant message
	session.AddAssistantMessage(resp.Content)

	// Return response
	c.JSON(http.StatusOK, ChatResponse{
		Response:  resp.Content,
		ToolCalls: resp.ToolCalls,
	})
}

// handleGetChatHistory handles requests to get chat history
func (s *Server) handleGetChatHistory(c *gin.Context) {
	// In a real implementation, this would retrieve the chat history from the session
	c.JSON(http.StatusOK, gin.H{"messages": []Message{}})
}

// ToolRequest represents a tool execution request
type ToolRequest struct {
	SessionID string                 `json:"sessionId"`
	Input     map[string]interface{} `json:"input"`
}

// handleExecuteTool handles requests to execute a tool
func (s *Server) handleExecuteTool(c *gin.Context) {
	toolName := c.Param("toolName")

	var req ToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get tool
	tool := s.toolRegistry.GetTool(toolName)
	if tool == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tool not found"})
		return
	}

	// Validate input
	if err := tool.ValidateInput(req.Input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute tool
	ctx := context.Background()
	output, err := tool.Execute(ctx, req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": output})
}

// handleListFiles handles requests to list files
func (s *Server) handleListFiles(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "."
	}

	// Use LSTool to list files
	lsTool := s.toolRegistry.GetTool("LS")
	if lsTool == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LS tool not found"})
		return
	}

	ctx := context.Background()
	output, err := lsTool.Execute(ctx, map[string]interface{}{
		"path": path,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// handleGetFile handles requests to get file content
func (s *Server) handleGetFile(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path is required"})
		return
	}

	// Use FileReadTool to read file
	fileReadTool := s.toolRegistry.GetTool("FileRead")
	if fileReadTool == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "FileRead tool not found"})
		return
	}

	ctx := context.Background()
	output, err := fileReadTool.Execute(ctx, map[string]interface{}{
		"file_path": path,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// FileUpdateRequest represents a file update request
type FileUpdateRequest struct {
	Content string `json:"content"`
}

// handleUpdateFile handles requests to update file content
func (s *Server) handleUpdateFile(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path is required"})
		return
	}

	var req FileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use FileWriteTool to write file
	fileWriteTool := s.toolRegistry.GetTool("FileWrite")
	if fileWriteTool == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "FileWrite tool not found"})
		return
	}

	ctx := context.Background()
	output, err := fileWriteTool.Execute(ctx, map[string]interface{}{
		"file_path": path,
		"content":   req.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// handleGetConfig handles requests to get configuration
func (s *Server) handleGetConfig(c *gin.Context) {
	globalConfig, err := s.config.GetGlobalConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get config"})
		return
	}

	// Don't expose API keys directly
	safeConfig := map[string]interface{}{
		"defaultModel":           globalConfig.DefaultModel,
		"hasCompletedOnboarding": globalConfig.HasCompletedOnboarding,
		"lastOnboardingVersion":  globalConfig.LastOnboardingVersion,
		"hasApiKeys": map[string]bool{
			"openai": globalConfig.APIKeys["openai"] != "",
		},
	}

	c.JSON(http.StatusOK, safeConfig)
}

// ConfigUpdateRequest represents a config update request
type ConfigUpdateRequest struct {
	DefaultModel string            `json:"defaultModel"`
	APIKeys      map[string]string `json:"apiKeys"`
}

// handleUpdateConfig handles requests to update configuration
func (s *Server) handleUpdateConfig(c *gin.Context) {
	var req ConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	globalConfig, err := s.config.GetGlobalConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get config"})
		return
	}

	// Update config
	if req.DefaultModel != "" {
		globalConfig.DefaultModel = req.DefaultModel
	}

	// Update API keys
	for provider, key := range req.APIKeys {
		if key != "" {
			globalConfig.APIKeys[provider] = key
		}
	}

	// Save config
	if err := s.config.SaveGlobalConfig(globalConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Run server
	if err := server.router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
