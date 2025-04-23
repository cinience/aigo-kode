package tools

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/cinience/aigo-kode/internal/core"
)

// GlobTool implements the Tool interface for finding files matching a pattern
type GlobTool struct{}

// Name returns the tool name
func (t *GlobTool) Name() string {
	return "Glob"
}

// Description returns the tool description
func (t *GlobTool) Description() string {
	return "Finds files matching a pattern"
}

// GlobToolOutput defines the output structure for GlobTool
type GlobToolOutput struct {
	Files []string `json:"files"`
	Error string   `json:"error,omitempty"`
}

// Execute executes the glob operation
func (t *GlobTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	// Extract pattern
	pattern, ok := input["pattern"].(string)
	if !ok || pattern == "" {
		return nil, errors.New("pattern is required and must be a string")
	}

	// Extract base directory (optional)
	baseDir := "."
	if baseDirVal, ok := input["base_dir"].(string); ok && baseDirVal != "" {
		baseDir = baseDirVal
	}

	// Ensure base directory exists
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return &GlobToolOutput{
			Files: []string{},
			Error: "Base directory does not exist",
		}, nil
	}

	// If pattern doesn't contain a path separator, search in all subdirectories
	if !strings.Contains(pattern, string(filepath.Separator)) {
		pattern = filepath.Join(baseDir, "**", pattern)
	} else if !filepath.IsAbs(pattern) {
		// If pattern is relative, make it relative to baseDir
		pattern = filepath.Join(baseDir, pattern)
	}

	// Find files matching pattern
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return &GlobToolOutput{
			Files: []string{},
			Error: "Invalid pattern: " + err.Error(),
		}, nil
	}

	// Filter out directories
	files := make([]string, 0, len(matches))
	for _, match := range matches {
		info, err := os.Stat(match)
		if err == nil && !info.IsDir() {
			files = append(files, match)
		}
	}

	return &GlobToolOutput{
		Files: files,
	}, nil
}

// ValidateInput validates the input parameters
func (t *GlobTool) ValidateInput(input map[string]interface{}) error {
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

	// Validate base_dir if present
	if baseDirVal, ok := input["base_dir"]; ok {
		baseDir, ok := baseDirVal.(string)
		if !ok {
			return errors.New("base_dir must be a string")
		}

		if baseDir != "" {
			if _, err := os.Stat(baseDir); os.IsNotExist(err) {
				return errors.New("base_dir does not exist")
			}
		}
	}

	return nil
}

// IsReadOnly returns whether the tool is read-only
func (t *GlobTool) IsReadOnly() bool {
	return true
}

// RequiresPermission checks if permission is needed
func (t *GlobTool) RequiresPermission(input map[string]interface{}) bool {
	return true
}

// NewGlobTool creates a new GlobTool
func NewGlobTool() core.Tool {
	return &GlobTool{}
}
