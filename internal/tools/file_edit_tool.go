package tools

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/cinience/aigo-kode/internal/core"
)

// FileEditTool implements the Tool interface for editing files
type FileEditTool struct{}

// Name returns the tool name
func (t *FileEditTool) Name() string {
	return "FileEdit"
}

// Description returns the tool description
func (t *FileEditTool) Description() string {
	return "Edits a file by replacing text"
}

// FileEditToolOutput defines the output structure for FileEditTool
type FileEditToolOutput struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Execute executes the file edit operation
func (t *FileEditTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	// Extract parameters
	filePath, ok := input["file_path"].(string)
	if !ok || filePath == "" {
		return nil, errors.New("file_path is required and must be a string")
	}

	oldText, ok := input["old_text"].(string)
	if !ok || oldText == "" {
		return nil, errors.New("old_text is required and must be a string")
	}

	newText, ok := input["new_text"].(string)
	if !ok {
		return nil, errors.New("new_text is required and must be a string")
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return &FileEditToolOutput{
			Success: false,
			Error:   "Failed to read file: " + err.Error(),
		}, nil
	}

	// Replace text
	newContent := strings.Replace(string(content), oldText, newText, -1)
	if newContent == string(content) {
		return &FileEditToolOutput{
			Success: false,
			Error:   "Old text not found in file",
		}, nil
	}

	// Write file
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return &FileEditToolOutput{
			Success: false,
			Error:   "Failed to write file: " + err.Error(),
		}, nil
	}

	return &FileEditToolOutput{
		Success: true,
	}, nil
}

// ValidateInput validates the input parameters
func (t *FileEditTool) ValidateInput(input map[string]interface{}) error {
	// Check if file_path exists and is a string
	filePathVal, ok := input["file_path"]
	if !ok {
		return errors.New("file_path is required")
	}

	filePath, ok := filePathVal.(string)
	if !ok {
		return errors.New("file_path must be a string")
	}

	if filePath == "" {
		return errors.New("file_path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	// Check if old_text exists and is a string
	oldTextVal, ok := input["old_text"]
	if !ok {
		return errors.New("old_text is required")
	}

	oldText, ok := oldTextVal.(string)
	if !ok {
		return errors.New("old_text must be a string")
	}

	if oldText == "" {
		return errors.New("old_text cannot be empty")
	}

	// Check if new_text exists and is a string
	newTextVal, ok := input["new_text"]
	if !ok {
		return errors.New("new_text is required")
	}

	_, ok = newTextVal.(string)
	if !ok {
		return errors.New("new_text must be a string")
	}

	return nil
}

func (t *FileEditTool) Arguments() string {
	return `{
		"file_path": {
			"type": "string",
			"description": "The path to the file to edit"
		},
		"old_text": {
			"type": "string",
			"description": "The text to replace"
		},
		"new_text": {
			"type": "string",
			"description": "The new text to replace with"
		}
	}`
}

// IsReadOnly returns whether the tool is read-only
func (t *FileEditTool) IsReadOnly() bool {
	return false
}

// RequiresPermission checks if permission is needed
func (t *FileEditTool) RequiresPermission(input map[string]interface{}) bool {
	return true
}

// NewFileEditTool creates a new FileEditTool
func NewFileEditTool() core.Tool {
	return &FileEditTool{}
}
