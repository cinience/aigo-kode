package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// GlobalConfig represents the global application configuration
type GlobalConfig struct {
	// DefaultModel is the default AI model to use
	DefaultModel string `json:"defaultModel"`
	// APIKeys maps provider names to API keys
	APIKeys map[string]string `json:"apiKeys"`
	// BaseURL is the base URL for the API
	BaseURL string `json:"baseURL"`
	// HasCompletedOnboarding indicates if the user has completed onboarding
	HasCompletedOnboarding bool `json:"hasCompletedOnboarding"`
	// LastOnboardingVersion is the version when onboarding was last completed
	LastOnboardingVersion string `json:"lastOnboardingVersion"`
}

// ProjectConfig represents project-specific configuration
type ProjectConfig struct {
	// ApprovedTools is a list of tools approved for use in this project
	ApprovedTools []string `json:"approvedTools"`
}

// Config defines the interface for configuration management
type Config interface {
	// GetGlobalConfig retrieves the global configuration
	GetGlobalConfig() (*GlobalConfig, error)

	// GetProjectConfig retrieves the project configuration
	GetProjectConfig(projectPath string) (*ProjectConfig, error)

	// SaveGlobalConfig saves the global configuration
	SaveGlobalConfig(config *GlobalConfig) error

	// SaveProjectConfig saves the project configuration
	SaveProjectConfig(projectPath string, config *ProjectConfig) error
}

// FileConfig implements Config using the filesystem
type FileConfig struct {
	// configDir is the directory where configuration files are stored
	configDir string
}

// NewFileConfig creates a new FileConfig
func NewFileConfig(configDir string) (*FileConfig, error) {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	return &FileConfig{
		configDir: configDir,
	}, nil
}

// GetGlobalConfig retrieves the global configuration
func (c *FileConfig) GetGlobalConfig() (*GlobalConfig, error) {
	configPath := filepath.Join(c.configDir, "config.json")

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return &GlobalConfig{
			DefaultModel: "gpt-3.5-turbo",
			APIKeys:      make(map[string]string),
		}, nil
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var config GlobalConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetProjectConfig retrieves the project configuration
func (c *FileConfig) GetProjectConfig(projectPath string) (*ProjectConfig, error) {
	if projectPath == "" {
		return nil, errors.New("project path cannot be empty")
	}

	configPath := filepath.Join(projectPath, ".aigo-kode.json")

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return &ProjectConfig{
			ApprovedTools: []string{},
		}, nil
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var config ProjectConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveGlobalConfig saves the global configuration
func (c *FileConfig) SaveGlobalConfig(config *GlobalConfig) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	configPath := filepath.Join(c.configDir, "config.json")
	return os.WriteFile(configPath, data, 0644)
}

// SaveProjectConfig saves the project configuration
func (c *FileConfig) SaveProjectConfig(projectPath string, config *ProjectConfig) error {
	if projectPath == "" {
		return errors.New("project path cannot be empty")
	}

	if config == nil {
		return errors.New("config cannot be nil")
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	configPath := filepath.Join(projectPath, ".aigo-kode.json")
	return os.WriteFile(configPath, data, 0644)
}
