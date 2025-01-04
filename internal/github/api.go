package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v57/github"
	"github.com/saint/ghquick/internal/log"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	logger *log.Logger
}

func NewClient(token string, debug bool) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return &Client{
		client: github.NewClient(tc),
		logger: log.New(debug),
	}
}

func (c *Client) EnsureRepositoryExists(ctx context.Context, name string, private bool) error {
	c.logger.Step("Checking if repository exists...")

	// Try to get the repository first
	_, _, err := c.client.Repositories.Get(ctx, os.Getenv("GITHUB_USERNAME"), name)
	if err == nil {
		c.logger.Info("Repository already exists")
		return nil
	}

	// Create repository if it doesn't exist
	c.logger.Step("Creating new repository: %s", name)
	repo := &github.Repository{
		Name:     github.String(name),
		Private:  github.Bool(private),
		AutoInit: github.Bool(false),
	}

	_, _, err = c.client.Repositories.Create(ctx, "", repo)
	if err != nil {
		c.logger.Error("Failed to create repository")
		return fmt.Errorf("failed to create repository: %w", err)
	}

	c.logger.Success("Repository created successfully")
	return nil
}
