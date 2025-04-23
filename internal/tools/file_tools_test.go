package tools

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileReadTool(t *testing.T) {
	// Create a temporary test file
	content := "This is a test file\nWith multiple lines\nFor testing FileReadTool"
	tmpFile, err := os.CreateTemp("", "filereadtest-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Create the tool
	tool := NewFileReadTool()

	// Test validation
	err = tool.ValidateInput(map[string]interface{}{
		"file_path": tmpFile.Name(),
	})
	assert.NoError(t, err)

	// Test validation with missing file_path
	err = tool.ValidateInput(map[string]interface{}{})
	assert.Error(t, err)

	// Test execution
	result, err := tool.Execute(context.Background(), map[string]interface{}{
		"file_path": tmpFile.Name(),
	})
	assert.NoError(t, err)

	// Check result
	fileResult, ok := result.(*FileReadToolOutput)
	assert.True(t, ok)
	assert.Equal(t, "text", fileResult.Type)
	assert.Equal(t, content, fileResult.Content)

	// Test with non-existent file
	result, err = tool.Execute(context.Background(), map[string]interface{}{
		"file_path": "non-existent-file.txt",
	})
	assert.NoError(t, err)
	fileResult, ok = result.(*FileReadToolOutput)
	assert.True(t, ok)
	assert.Equal(t, "error", fileResult.Type)
	assert.Contains(t, fileResult.Error, "does not exist")
}

func TestFileWriteTool(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "filewritetest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFilePath := tmpDir + "/test.txt"

	// Create the tool
	tool := NewFileWriteTool()

	// Test validation
	err = tool.ValidateInput(map[string]interface{}{
		"file_path": testFilePath,
		"content":   "Test content",
	})
	assert.NoError(t, err)

	// Test validation with missing parameters
	err = tool.ValidateInput(map[string]interface{}{
		"file_path": testFilePath,
	})
	assert.Error(t, err)

	// Test execution - write file
	result, err := tool.Execute(context.Background(), map[string]interface{}{
		"file_path": testFilePath,
		"content":   "Test content",
	})
	assert.NoError(t, err)

	// Check result
	writeResult, ok := result.(*FileWriteToolOutput)
	assert.True(t, ok)
	assert.True(t, writeResult.Success)

	// Verify file was written
	content, err := os.ReadFile(testFilePath)
	assert.NoError(t, err)
	assert.Equal(t, "Test content", string(content))

	// Test append mode
	result, err = tool.Execute(context.Background(), map[string]interface{}{
		"file_path": testFilePath,
		"content":   "\nAppended content",
		"append":    true,
	})
	assert.NoError(t, err)

	// Verify content was appended
	content, err = os.ReadFile(testFilePath)
	assert.NoError(t, err)
	assert.Equal(t, "Test content\nAppended content", string(content))
}
