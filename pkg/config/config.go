package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AgentAdapter string `yaml:"agent_adapter"`
	Atlassian    struct {
		BaseURL string `yaml:"base_url"`
		User    string `yaml:"user"`
		APIKey  string `yaml:"api_key"`
	} `yaml:"atlassian"`
	MaxConcurrency int `yaml:"max_concurrency"`
	MaxRetries     int `yaml:"max_retries"`

	// LLM is the internal configuration state populated from agent_adapter
	LLM struct {
		Provider string
		Model    string
		APIKey   string
	}
}

// DefaultPath returns the default location for the config file: ~/.github.com/jefflunt/breakdown/config.yml
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".github.com/jefflunt/breakdown/config.yml" // fallback
	}
	return filepath.Join(home, ".breakdown", "config.yml")
}

// expandTilde handles resolving the ~ symbol in file paths
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func LoadConfig(path string) (*Config, error) {
	expandedPath := expandTilde(path)
	data, err := os.ReadFile(expandedPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if no file exists
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file at %s: %w", expandedPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file at %s: %w", expandedPath, err)
	}

	// Load Atlassian config from environment variables if not set in config file
	if cfg.Atlassian.BaseURL == "" {
		cfg.Atlassian.BaseURL = os.Getenv("BREAKDOWN_ATLASSIAN_BASE_URL")
	}
	if cfg.Atlassian.User == "" {
		cfg.Atlassian.User = os.Getenv("BREAKDOWN_ATLASSIAN_API_USER")
	}
	if cfg.Atlassian.APIKey == "" {
		cfg.Atlassian.APIKey = os.Getenv("BREAKDOWN_ATLASSIAN_API_KEY")
	}

	// Load Agent Adapter config from environment variables if not set in config file
	if cfg.AgentAdapter == "" {
		cfg.AgentAdapter = os.Getenv("BREAKDOWN_AGENT_ADAPTER")
	}
	if cfg.AgentAdapter == "" {
		cfg.AgentAdapter = "gemini:google/gemini-3.1-flash-lite-preview"
	}

	// Parse the agent adapter and populate internal LLM config
	cliName, _, model, err := ParseAdapter(cfg.AgentAdapter)
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent adapter: %w", err)
	}
	cfg.LLM.Provider = cliName
	cfg.LLM.Model = model

	return &cfg, nil
}

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.AgentAdapter = "gemini:google/gemini-3.1-flash-lite-preview"
	cfg.LLM.Provider = "gemini"
	cfg.LLM.Model = "gemini-3.1-flash-lite-preview"
	return cfg
}

// ParseAdapter decomposes the cliName:provider/model string format.
// It returns the CLI tool name, the LLM provider, the specific model, and any parsing error.
func ParseAdapter(adapter string) (string, string, string, error) {
	if adapter == "" {
		return "", "", "", fmt.Errorf("empty adapter string")
	}

	parts := strings.SplitN(adapter, ":", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid adapter format: %q (expected 'cliName:provider/model')", adapter)
	}
	cliName := strings.TrimSpace(parts[0])

	subParts := strings.SplitN(parts[1], "/", 2)
	if len(subParts) != 2 {
		return "", "", "", fmt.Errorf("invalid adapter format: %q (expected 'cliName:provider/model')", adapter)
	}
	provider := strings.TrimSpace(subParts[0])
	model := strings.TrimSpace(subParts[1])

	if cliName == "" || provider == "" || model == "" {
		return "", "", "", fmt.Errorf("invalid adapter: components cannot be empty")
	}

	return cliName, provider, model, nil
}
