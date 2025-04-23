package tools

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashTool(t *testing.T) {
	// Create the tool
	tool := NewBashTool()

	// Test validation
	err := tool.ValidateInput(map[string]interface{}{
		"command": "echo 'Hello World'",
	})
	assert.NoError(t, err)

	// Test validation with missing command
	err = tool.ValidateInput(map[string]interface{}{})
	assert.Error(t, err)

	// Test validation with banned command
	err = tool.ValidateInput(map[string]interface{}{
		"command": "rm -rf /",
	})
	assert.Error(t, err)

	// Test execution with simple command
	result, err := tool.Execute(context.Background(), map[string]interface{}{
		"command": "echo 'Hello World'",
	})
	assert.NoError(t, err)

	// Check result
	bashResult, ok := result.(*BashToolOutput)
	assert.True(t, ok)
	assert.Contains(t, bashResult.Stdout, "Hello World")
	assert.False(t, bashResult.Interrupted)

	// Test with command that produces stderr
	result, err = tool.Execute(context.Background(), map[string]interface{}{
		"command": "ls /nonexistent",
	})
	assert.NoError(t, err)

	// Check result has stderr
	bashResult, ok = result.(*BashToolOutput)
	assert.True(t, ok)
	assert.NotEmpty(t, bashResult.Stderr)
}

func TestGlobTool(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "globtest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some test files
	testFiles := []string{"test1.txt", "test2.txt", "test.go", "subdir/test3.txt"}
	for _, file := range testFiles {
		filePath := tmpDir + "/" + file
		dirPath := tmpDir + "/" + file[:len(file)-len(file)]
		os.MkdirAll(dirPath, 0755)
		os.WriteFile(filePath, []byte("test content"), 0644)
	}

	// Create the tool
	tool := NewGlobTool()

	// Test validation
	err = tool.ValidateInput(map[string]interface{}{
		"pattern":  "*.txt",
		"base_dir": tmpDir,
	})
	assert.NoError(t, err)

	// Test validation with missing pattern
	err = tool.ValidateInput(map[string]interface{}{
		"base_dir": tmpDir,
	})
	assert.Error(t, err)

	// Test execution
	result, err := tool.Execute(context.Background(), map[string]interface{}{
		"pattern":  "*.txt",
		"base_dir": tmpDir,
	})
	assert.NoError(t, err)

	// Check result
	globResult, ok := result.(*GlobToolOutput)
	assert.True(t, ok)
	assert.Len(t, globResult.Files, 2) // Should find test1.txt and test2.txt

	// Test with more specific pattern
	result, err = tool.Execute(context.Background(), map[string]interface{}{
		"pattern":  "test1.txt",
		"base_dir": tmpDir,
	})
	assert.NoError(t, err)

	// Check result
	globResult, ok = result.(*GlobToolOutput)
	assert.True(t, ok)
	assert.Len(t, globResult.Files, 1) // Should find only test1.txt
}
