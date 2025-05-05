package tools

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/cinience/aigo-kode/internal/core"
)

// FileWriteTool implements the Tool interface for writing to files
type FileWriteTool struct{}

// Name returns the tool name
func (t *FileWriteTool) Name() string {
	return "FileWrite"
}

// Description returns the tool description
func (t *FileWriteTool) Description() string {
	return "Writes content to a file"
}

// FileWriteToolOutput defines the output structure for FileWriteTool
type FileWriteToolOutput struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Execute executes the file write operation
func (t *FileWriteTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	// Extract file path
	filePath, ok := input["file_path"].(string)
	if !ok || filePath == "" {
		return nil, errors.New("file_path is required and must be a string")
	}

	// Extract content
	content, ok := input["content"].(string)
	if !ok {
		return nil, errors.New("content is required and must be a string")
	}

	// Extract optional parameters
	append := false
	if appendVal, ok := input["append"].(bool); ok {
		append = appendVal
	}

	leadingNewline := append
	if leadingNewlineVal, ok := input["leading_newline"].(bool); ok {
		leadingNewline = leadingNewlineVal
	}

	trailingNewline := true
	if trailingNewlineVal, ok := input["trailing_newline"].(bool); ok {
		trailingNewline = trailingNewlineVal
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &FileWriteToolOutput{
			Success: false,
			Error:   "Failed to create directory: " + err.Error(),
		}, nil
	}

	// Prepare content with newlines if needed
	if leadingNewline {
		content = "\n" + content
	}
	if trailingNewline && !append {
		content = content + "\n"
	}

	var err error
	if append {
		// Open file in append mode
		var file *os.File
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			_, err = file.WriteString(content)
			file.Close()
		}
	} else {
		// Write or overwrite file
		err = os.WriteFile(filePath, []byte(content), 0644)
	}

	if err != nil {
		return &FileWriteToolOutput{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &FileWriteToolOutput{
		Success: true,
	}, nil
}

// ValidateInput validates the input parameters
func (t *FileWriteTool) ValidateInput(input map[string]interface{}) error {
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

	// Check if content exists
	contentVal, ok := input["content"]
	if !ok {
		return errors.New("content is required")
	}

	_, ok = contentVal.(string)
	if !ok {
		return errors.New("content must be a string")
	}

	// Validate append if present
	if appendVal, ok := input["append"]; ok {
		_, ok := appendVal.(bool)
		if !ok {
			return errors.New("append must be a boolean")
		}
	}

	// Validate leading_newline if present
	if leadingNewlineVal, ok := input["leading_newline"]; ok {
		_, ok := leadingNewlineVal.(bool)
		if !ok {
			return errors.New("leading_newline must be a boolean")
		}
	}

	// Validate trailing_newline if present
	if trailingNewlineVal, ok := input["trailing_newline"]; ok {
		_, ok := trailingNewlineVal.(bool)
		if !ok {
			return errors.New("trailing_newline must be a boolean")
		}
	}

	return nil
}

func (t *FileWriteTool) Arguments() string {
	return `{
		"file_path": {
			"type": "string",
			"description": "The path of the file to write to"
		},
		"content": {
			"type": "string",
			"description": "The content to write to the file"
		},
		"append": {
			"type": "boolean",
			"description": "Whether to append to the file or overwrite it (default: false)"
		},
		"leading_newline": {
			"type": "boolean",
			"description": "Whether to add a leading newline to the content"
		},
		"trailing_newline": {
			"type": "boolean",
			"description": "Whether to add a trailing newline to the content (default: true)"
		}
	}`
}


// IsReadOnly returns whether the tool is read-only
func (t *FileWriteTool) IsReadOnly() bool {
	return false
}

// RequiresPermission checks if permission is needed
func (t *FileWriteTool) RequiresPermission(input map[string]interface{}) bool {
	return true
}

// NewFileWriteTool creates a new FileWriteTool
func NewFileWriteTool() core.Tool {
	return &FileWriteTool{}
}
