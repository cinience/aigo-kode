package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileConfig(t *testing.T) {
	// Create a temporary directory for config files
	tmpDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a new FileConfig
	config, err := NewFileConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Test GetGlobalConfig with no existing config file
	globalConfig, err := config.GetGlobalConfig()
	assert.NoError(t, err)
	assert.NotNil(t, globalConfig)
	assert.Equal(t, "gpt-3.5-turbo", globalConfig.DefaultModel)
	assert.NotNil(t, globalConfig.APIKeys)

	// Test SaveGlobalConfig
	globalConfig.DefaultModel = "gpt-4"
	globalConfig.APIKeys["openai"] = "test-api-key"
	globalConfig.HasCompletedOnboarding = true
	globalConfig.LastOnboardingVersion = "1.0.0"

	err = config.SaveGlobalConfig(globalConfig)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(filepath.Join(tmpDir, "config.json"))
	assert.NoError(t, err)

	// Test GetGlobalConfig with existing config file
	loadedConfig, err := config.GetGlobalConfig()
	assert.NoError(t, err)
	assert.Equal(t, "gpt-4", loadedConfig.DefaultModel)
	assert.Equal(t, "test-api-key", loadedConfig.APIKeys["openai"])
	assert.True(t, loadedConfig.HasCompletedOnboarding)
	assert.Equal(t, "1.0.0", loadedConfig.LastOnboardingVersion)

	// Test project config
	projectDir, err := os.MkdirTemp("", "projecttest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(projectDir)

	// Test GetProjectConfig with no existing config file
	projectConfig, err := config.GetProjectConfig(projectDir)
	assert.NoError(t, err)
	assert.NotNil(t, projectConfig)
	assert.Empty(t, projectConfig.ApprovedTools)

	// Test SaveProjectConfig
	projectConfig.ApprovedTools = []string{"BashTool", "FileReadTool"}
	err = config.SaveProjectConfig(projectDir, projectConfig)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(filepath.Join(projectDir, ".go-anon-kode.json"))
	assert.NoError(t, err)

	// Test GetProjectConfig with existing config file
	loadedProjectConfig, err := config.GetProjectConfig(projectDir)
	assert.NoError(t, err)
	assert.Equal(t, []string{"BashTool", "FileReadTool"}, loadedProjectConfig.ApprovedTools)

	// Test error cases
	err = config.SaveGlobalConfig(nil)
	assert.Error(t, err)

	err = config.SaveProjectConfig("", projectConfig)
	assert.Error(t, err)

	err = config.SaveProjectConfig(projectDir, nil)
	assert.Error(t, err)
}
