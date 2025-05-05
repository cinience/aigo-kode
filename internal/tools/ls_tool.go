package tools

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/cinience/aigo-kode/internal/core"
)

// LSTool implements the Tool interface for listing directory contents
type LSTool struct{}

// Name returns the tool name
func (t *LSTool) Name() string {
	return "LS"
}

// Description returns the tool description
func (t *LSTool) Description() string {
	return "Lists files and directories in a specified path"
}

// LSEntry represents a file or directory entry
type LSEntry struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	IsDir     bool   `json:"is_dir"`
	Size      int64  `json:"size"`
	Extension string `json:"extension,omitempty"`
}

// LSToolOutput defines the output structure for LSTool
type LSToolOutput struct {
	Entries []LSEntry `json:"entries"`
	Error   string    `json:"error,omitempty"`
}

// Execute executes the ls operation
func (t *LSTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	// Extract directory path
	dirPath, ok := input["path"].(string)
	if !ok {
		dirPath = "." // Default to current directory
	}

	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		return &LSToolOutput{
			Entries: []LSEntry{},
			Error:   "Directory does not exist or cannot be accessed: " + err.Error(),
		}, nil
	}

	// If path is a file, return info about just that file
	if !info.IsDir() {
		entry := LSEntry{
			Name:      filepath.Base(dirPath),
			Path:      dirPath,
			IsDir:     false,
			Size:      info.Size(),
			Extension: filepath.Ext(dirPath),
		}
		return &LSToolOutput{
			Entries: []LSEntry{entry},
		}, nil
	}

	// Read directory contents
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return &LSToolOutput{
			Entries: []LSEntry{},
			Error:   "Failed to read directory: " + err.Error(),
		}, nil
	}

	// Extract optional parameters
	showHidden := false
	if showHiddenVal, ok := input["show_hidden"].(bool); ok {
		showHidden = showHiddenVal
	}

	sortBy := "name"
	if sortByVal, ok := input["sort_by"].(string); ok {
		sortBy = sortByVal
	}

	// Process directory entries
	entries := make([]LSEntry, 0, len(files))
	for _, file := range files {
		// Skip hidden files if not showing them
		if !showHidden && file.Name()[0] == '.' {
			continue
		}

		entry := LSEntry{
			Name:  file.Name(),
			Path:  filepath.Join(dirPath, file.Name()),
			IsDir: file.IsDir(),
			Size:  file.Size(),
		}

		if !file.IsDir() {
			entry.Extension = filepath.Ext(file.Name())
		}

		entries = append(entries, entry)
	}

	// Sort entries
	switch sortBy {
	case "name":
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name < entries[j].Name
		})
	case "size":
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Size < entries[j].Size
		})
	case "type":
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].IsDir && !entries[j].IsDir {
				return true
			}
			if !entries[i].IsDir && entries[j].IsDir {
				return false
			}
			return entries[i].Name < entries[j].Name
		})
	}

	return &LSToolOutput{
		Entries: entries,
	}, nil
}

// ValidateInput validates the input parameters
func (t *LSTool) ValidateInput(input map[string]interface{}) error {
	// path is optional, but if provided must be a string
	if pathVal, ok := input["path"]; ok {
		path, ok := pathVal.(string)
		if !ok {
			return errors.New("path must be a string")
		}

		if path != "" {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return errors.New("path does not exist")
			}
		}
	}

	// Validate show_hidden if present
	if showHiddenVal, ok := input["show_hidden"]; ok {
		_, ok := showHiddenVal.(bool)
		if !ok {
			return errors.New("show_hidden must be a boolean")
		}
	}

	// Validate sort_by if present
	if sortByVal, ok := input["sort_by"]; ok {
		sortBy, ok := sortByVal.(string)
		if !ok {
			return errors.New("sort_by must be a string")
		}

		validSortOptions := map[string]bool{
			"name": true,
			"size": true,
			"type": true,
		}

		if !validSortOptions[sortBy] {
			return errors.New("sort_by must be one of: name, size, type")
		}
	}

	return nil
}

func (t *LSTool) Arguments() string {
	return `{
		"path": {
			"type": "string",
			"description": "The path to the directory to list (optional)"
		},
		"show_hidden": {
			"type": "boolean",
			"description": "Whether to show hidden files and directories (optional)"
		},
		"sort_by": {
			"type": "string",
			"description": "The field to sort by (optional): name, size, type"
		}
		`
}

// IsReadOnly returns whether the tool is read-only
func (t *LSTool) IsReadOnly() bool {
	return true
}

// RequiresPermission checks if permission is needed
func (t *LSTool) RequiresPermission(input map[string]interface{}) bool {
	return true
}

// NewLSTool creates a new LSTool
func NewLSTool() core.Tool {
	return &LSTool{}
}
