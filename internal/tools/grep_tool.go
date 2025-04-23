package tools

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/cinience/aigo-kode/internal/core"
)

// GrepTool implements the Tool interface for searching file contents
type GrepTool struct{}

// Name returns the tool name
func (t *GrepTool) Name() string {
	return "Grep"
}

// Description returns the tool description
func (t *GrepTool) Description() string {
	return "Searches for text patterns in files"
}

// GrepToolOutput defines the output structure for GrepTool
type GrepMatch struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Content string `json:"content"`
}

type GrepToolOutput struct {
	Matches []GrepMatch `json:"matches"`
	Error   string      `json:"error,omitempty"`
}

// Execute executes the grep operation
func (t *GrepTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	// Extract pattern
	pattern, ok := input["pattern"].(string)
	if !ok || pattern == "" {
		return nil, errors.New("pattern is required and must be a string")
	}

	// Extract file paths or patterns
	var filePaths []string
	if filePathsVal, ok := input["file_paths"].([]interface{}); ok {
		for _, pathVal := range filePathsVal {
			if path, ok := pathVal.(string); ok && path != "" {
				filePaths = append(filePaths, path)
			}
		}
	} else if filePathVal, ok := input["file_paths"].(string); ok && filePathVal != "" {
		filePaths = []string{filePathVal}
	}

	if len(filePaths) == 0 {
		return nil, errors.New("file_paths is required and must be a string or array of strings")
	}

	// Extract max matches (optional)
	maxMatches := 100
	if maxMatchesVal, ok := input["max_matches"].(float64); ok {
		maxMatches = int(maxMatchesVal)
	}

	// Process each file path
	var allMatches []GrepMatch
	for _, path := range filePaths {
		// Handle glob patterns
		matches, err := filepath.Glob(path)
		if err != nil || len(matches) == 0 {
			// If not a glob pattern or no matches, treat as a single file
			matches = []string{path}
		}

		// Process each matched file
		for _, filePath := range matches {
			// Check if file exists and is a regular file
			info, err := os.Stat(filePath)
			if err != nil || info.IsDir() {
				continue
			}

			// Read file content
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			// Search for pattern in each line
			lines := strings.Split(string(content), "\n")
			for i, line := range lines {
				if strings.Contains(line, pattern) {
					allMatches = append(allMatches, GrepMatch{
						File:    filePath,
						Line:    i + 1,
						Content: line,
					})

					// Check if we've reached the maximum number of matches
					if len(allMatches) >= maxMatches {
						break
					}
				}
			}

			// Check if we've reached the maximum number of matches
			if len(allMatches) >= maxMatches {
				break
			}
		}

		// Check if we've reached the maximum number of matches
		if len(allMatches) >= maxMatches {
			break
		}
	}

	return &GrepToolOutput{
		Matches: allMatches,
	}, nil
}

// ValidateInput validates the input parameters
func (t *GrepTool) ValidateInput(input map[string]interface{}) error {
	// Check if pattern exists and is a string
	patternVal, ok := input["pattern"]
	if !ok {
		return errors.New("pattern is required")
	}

	pattern, ok := patternVal.(string)
	if !ok {
		return errors.New("pattern must be a string")
	}

	if pattern == "" {
		return errors.New("pattern cannot be empty")
	}

	// Check if file_paths exists and is a string or array of strings
	filePathsVal, ok := input["file_paths"]
	if !ok {
		return errors.New("file_paths is required")
	}

	// Check if file_paths is a string
	if _, ok := filePathsVal.(string); !ok {
		// If not a string, check if it's an array
		filePathsArray, ok := filePathsVal.([]interface{})
		if !ok || len(filePathsArray) == 0 {
			return errors.New("file_paths must be a non-empty string or array of strings")
		}

		// Check that all elements in the array are strings
		for _, pathVal := range filePathsArray {
			if _, ok := pathVal.(string); !ok {
				return errors.New("all elements in file_paths array must be strings")
			}
		}
	}

	// Validate max_matches if present
	if maxMatchesVal, ok := input["max_matches"]; ok {
		maxMatches, ok := maxMatchesVal.(float64)
		if !ok {
			return errors.New("max_matches must be a number")
		}

		if maxMatches <= 0 {
			return errors.New("max_matches must be positive")
		}
	}

	return nil
}

// IsReadOnly returns whether the tool is read-only
func (t *GrepTool) IsReadOnly() bool {
	return true
}

// RequiresPermission checks if permission is needed
func (t *GrepTool) RequiresPermission(input map[string]interface{}) bool {
	return true
}

// NewGrepTool creates a new GrepTool
func NewGrepTool() core.Tool {
	return &GrepTool{}
}
