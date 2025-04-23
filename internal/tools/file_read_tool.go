package tools

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cinience/aigo-kode/internal/core"
)

// FileReadTool implements the Tool interface for reading file contents
type FileReadTool struct{}

// Name returns the tool name
func (t *FileReadTool) Name() string {
	return "FileRead"
}

// Description returns the tool description
func (t *FileReadTool) Description() string {
	return "Reads the contents of a file"
}

// FileReadToolOutput defines the output structure for FileReadTool
type FileReadToolOutput struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Execute executes the file read operation
func (t *FileReadTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	// Extract file path
	filePath, ok := input["file_path"].(string)
	if !ok || filePath == "" {
		return nil, errors.New("file_path is required and must be a string")
	}

	// Extract optional parameters
	var offset, limit int
	if offsetVal, ok := input["offset"].(float64); ok {
		offset = int(offsetVal)
	}
	if limitVal, ok := input["limit"].(float64); ok {
		limit = int(limitVal)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &FileReadToolOutput{
			Type:  "error",
			Error: "File does not exist",
		}, nil
	}

	// Check if it's an image file
	ext := strings.ToLower(filepath.Ext(filePath))
	imageExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".bmp": true, ".webp": true,
	}

	if imageExts[ext] {
		return &FileReadToolOutput{
			Type: "image",
		}, nil
	}

	// Read file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &FileReadToolOutput{
			Type:  "error",
			Error: err.Error(),
		}, nil
	}

	// Convert to string
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	// Apply offset and limit if provided
	if offset > 0 && offset < len(lines) {
		lines = lines[offset:]
	}
	if limit > 0 && limit < len(lines) {
		lines = lines[:limit]
	}

	return &FileReadToolOutput{
		Type:    "text",
		Content: strings.Join(lines, "\n"),
	}, nil
}

// ValidateInput validates the input parameters
func (t *FileReadTool) ValidateInput(input map[string]interface{}) error {
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

	// Validate offset if present
	if offsetVal, ok := input["offset"]; ok {
		offset, ok := offsetVal.(float64)
		if !ok {
			return errors.New("offset must be a number")
		}

		if offset < 0 {
			return errors.New("offset must be non-negative")
		}
	}

	// Validate limit if present
	if limitVal, ok := input["limit"]; ok {
		limit, ok := limitVal.(float64)
		if !ok {
			return errors.New("limit must be a number")
		}

		if limit <= 0 {
			return errors.New("limit must be positive")
		}
	}

	return nil
}

// IsReadOnly returns whether the tool is read-only
func (t *FileReadTool) IsReadOnly() bool {
	return true
}

// RequiresPermission checks if permission is needed
func (t *FileReadTool) RequiresPermission(input map[string]interface{}) bool {
	return true
}

// NewFileReadTool creates a new FileReadTool
func NewFileReadTool() core.Tool {
	return &FileReadTool{}
}
