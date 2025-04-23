package tools

import (
	"github.com/cinience/aigo-kode/internal/core"
)

// ToolRegistry manages the collection of available tools
type ToolRegistry struct {
	tools map[string]func() core.Tool
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]func() core.Tool),
	}
}

// RegisterTool registers a tool factory function
func (r *ToolRegistry) RegisterTool(name string, factory func() core.Tool) {
	r.tools[name] = factory
}

// GetTool creates a new instance of a tool by name
func (r *ToolRegistry) GetTool(name string) core.Tool {
	if factory, ok := r.tools[name]; ok {
		return factory()
	}
	return nil
}

// GetAllTools returns all registered tools
func (r *ToolRegistry) GetAllTools() []core.Tool {
	tools := make([]core.Tool, 0, len(r.tools))
	for _, factory := range r.tools {
		tools = append(tools, factory())
	}
	return tools
}

// GetReadOnlyTools returns all read-only tools
func (r *ToolRegistry) GetReadOnlyTools() []core.Tool {
	tools := make([]core.Tool, 0, len(r.tools))
	for _, factory := range r.tools {
		tool := factory()
		if tool.IsReadOnly() {
			tools = append(tools, tool)
		}
	}
	return tools
}

// DefaultToolRegistry creates and returns a registry with all standard tools
func DefaultToolRegistry() *ToolRegistry {
	registry := NewToolRegistry()

	// Register all tools
	registry.RegisterTool("Bash", NewBashTool)
	registry.RegisterTool("FileRead", NewFileReadTool)
	registry.RegisterTool("FileWrite", NewFileWriteTool)
	registry.RegisterTool("FileEdit", NewFileEditTool)
	registry.RegisterTool("Glob", NewGlobTool)
	registry.RegisterTool("Grep", NewGrepTool)
	registry.RegisterTool("LS", NewLSTool)
	registry.RegisterTool("Think", NewThinkTool)

	return registry
}
