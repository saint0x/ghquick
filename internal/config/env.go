package config

import (
	"errors"
	"os"
)

const (
	EnvGitHubToken    = "GITHUB_TOKEN"
	EnvGitHubUsername = "GITHUB_USERNAME"
	EnvOpenAIKey      = "OPENAI_API_KEY"
)

type Config struct {
	GitHubToken    string
	GitHubUsername string
	OpenAIKey      string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	githubToken := os.Getenv(EnvGitHubToken)
	githubUsername := os.Getenv(EnvGitHubUsername)
	openAIKey := os.Getenv(EnvOpenAIKey)

	if githubToken == "" {
		return nil, errors.New("GITHUB_TOKEN environment variable is required")
	}
	if githubUsername == "" {
		return nil, errors.New("GITHUB_USERNAME environment variable is required")
	}
	if openAIKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is required")
	}

	return &Config{
		GitHubToken:    githubToken,
		GitHubUsername: githubUsername,
		OpenAIKey:      openAIKey,
	}, nil
}
