package github

import (
	"context"
	"fmt"
	"net/http"
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

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	if githubErr, ok := err.(*github.ErrorResponse); ok {
		return githubErr.Response.StatusCode == http.StatusNotFound
	}
	return false
}

func (c *Client) EnsureRepositoryExists(ctx context.Context, name string, private bool) error {
	c.logger.Step("Checking if repository exists...")
	username := os.Getenv("GITHUB_USERNAME")

	// Try to get the repository first
	repo, _, err := c.client.Repositories.Get(ctx, username, name)
	if err == nil {
		c.logger.Info("Repository exists, will append changes")
		// Update repository settings if needed
		if repo.GetPrivate() != private {
			c.logger.Step("Updating repository visibility...")
			repo.Private = github.Bool(private)
			_, _, err = c.client.Repositories.Edit(ctx, username, name, repo)
			if err != nil {
				c.logger.Error("Failed to update repository visibility")
				return fmt.Errorf("failed to update repository: %w", err)
			}
			c.logger.Success("Repository visibility updated")
		}
		return nil
	}

	// Only create if repository doesn't exist
	if isNotFound(err) {
		c.logger.Step("Repository doesn't exist, creating new repository: %s", name)
		repo = &github.Repository{
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

	// If we get here, it's an unexpected error
	c.logger.Error("Failed to check repository")
	return fmt.Errorf("failed to check repository: %w", err)
}
