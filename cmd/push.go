package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/saint/ghquick/internal/ai"
	"github.com/saint/ghquick/internal/cache"
	"github.com/saint/ghquick/internal/config"
	"github.com/saint/ghquick/internal/git"
	"github.com/saint/ghquick/internal/github"
	"github.com/saint/ghquick/internal/log"
	"github.com/spf13/cobra"
)

var (
	repoName   string
	commitMsg  string
	autoCommit bool
	repoCache  *cache.RepoCache
	debug      bool
	logger     *log.Logger
	private    bool
)

func init() {
	rootCmd.AddCommand(pushCmd)
	repoCache = cache.NewRepoCache()

	pushCmd.Flags().StringVar(&repoName, "name", "", "Repository name (defaults to current directory name)")
	pushCmd.Flags().StringVar(&commitMsg, "commitmsg", "", "Commit message")
	pushCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug logging")
	pushCmd.Flags().BoolVar(&private, "private", false, "Create repository as private")
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes to GitHub",
	Long: `Push changes to GitHub with optional AI-powered commit messages.
Example: 
  ghquick push start        # AI-powered push with automatic commit message
  ghquick push --name my-repo --commitmsg "feature: new stuff"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger = log.New(debug)
		if len(args) > 0 && args[0] == "start" {
			autoCommit = true
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Load configuration
		logger.Step("Loading configuration...")
		cfg, err := config.LoadFromEnv()
		if err != nil {
			logger.Error("Failed to load configuration")
			return fmt.Errorf("failed to load config: %w", err)
		}
		logger.Success("Configuration loaded")

		// Get current working directory
		wd, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get working directory")
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		// If repo name is not provided, use current directory name
		if repoName == "" {
			repoName = filepath.Base(wd)
			logger.Info("Using current directory name as repository name: %s", repoName)
		}

		// Initialize services
		gitOps := git.NewOperations(wd, debug)
		ghClient := github.NewClient(cfg.GitHubToken, debug)
		commitGen := ai.NewCommitMessageGenerator(cfg.OpenAIKey)

		// Ensure GitHub repository exists
		if err := ghClient.EnsureRepositoryExists(ctx, repoName, private); err != nil {
			return fmt.Errorf("failed to ensure repository exists: %w", err)
		}

		// Ensure git is set up
		if err := gitOps.EnsureGitSetup(ctx, repoName); err != nil {
			return fmt.Errorf("failed to setup git: %w", err)
		}

		// Stage all files first
		if err := gitOps.StageAll(ctx); err != nil {
			if err.Error() == "no changes to commit" {
				logger.Warning("No changes to commit")
				return nil
			}
			return fmt.Errorf("failed to stage files: %w", err)
		}

		// Get diff for commit message generation
		diff, err := gitOps.GetDiff(ctx)
		if err != nil {
			return fmt.Errorf("failed to get diff: %w", err)
		}

		// Generate commit message if needed
		if autoCommit {
			logger.Step("Generating commit message...")
			result := make(chan ai.GenerateResult, 1)
			commitGen.GenerateFromDiffAsync(ctx, diff, result)

			select {
			case res := <-result:
				if res.Error != nil {
					logger.Error("Failed to generate commit message")
					return fmt.Errorf("failed to generate commit message: %w", res.Error)
				}
				commitMsg = res.Message
				logger.Success("Commit message generated: %s", commitMsg)
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if commitMsg == "" {
			logger.Error("Commit message is required")
			return fmt.Errorf("commit message is required (use --commitmsg or 'start' for AI-generated message)")
		}

		// Commit changes
		if err := gitOps.Commit(ctx, commitMsg); err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}

		// Push changes
		if err := gitOps.Push(ctx, "origin", "main"); err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}

		logger.Success("🚀 Successfully pushed changes to GitHub!")
		return nil
	},
}
