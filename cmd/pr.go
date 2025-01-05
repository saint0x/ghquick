package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/saint0x/ghquick-cli/internal/ai"
	"github.com/saint0x/ghquick-cli/internal/config"
	"github.com/saint0x/ghquick-cli/internal/git"
	"github.com/saint0x/ghquick-cli/internal/github"
	"github.com/saint0x/ghquick-cli/internal/log"
	"github.com/spf13/cobra"
)

var (
	prTitle    string
	prBody     string
	prNumber   int
	baseBranch string
)

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.AddCommand(createPRCmd)
	prCmd.AddCommand(mergePRCmd)

	// Shared flags
	prCmd.PersistentFlags().StringVar(&repoName, "name", "", "Repository name (defaults to current directory name)")
	prCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	prCmd.PersistentFlags().DurationVar(&timeout, "timeout", timeout, "Timeout for operations (default 2m)")

	// Create PR specific flags
	createPRCmd.Flags().StringVar(&prTitle, "title", "", "Pull request title")
	createPRCmd.Flags().StringVar(&prBody, "body", "", "Pull request body")
	createPRCmd.Flags().StringVar(&baseBranch, "base", "main", "Base branch to create PR against")
	createPRCmd.Flags().StringVar(&commitMsg, "commitmsg", "", "Commit message before creating PR")

	// Merge PR specific flags
	mergePRCmd.Flags().IntVar(&prNumber, "number", 0, "Pull request number to merge")
}

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
	Long: `Create and manage pull requests with ease.

Available Actions:
  • create: Create a new PR with optional AI-generated commit message
  • merge:  Merge PRs with interactive selection or by number

Examples:
  ghquick pr create start              # Create PR with AI commit message
  ghquick pr create --title "Fix bug"  # Create PR with custom title
  ghquick pr merge                     # Interactive PR selection
  ghquick pr merge --number 123        # Merge specific PR`,
}

var createPRCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pull request",
	Long: `Create a pull request from a new branch to the target branch.
The command will:
1. Create a new branch from target branch
2. Stage and commit changes
3. Create PR with AI-generated title (or custom via flags)

Examples:
  ghquick pr create                    # Create PR with AI-generated title
  ghquick pr create --title "fix bug"  # Create PR with custom title
  ghquick pr create --base develop     # Create PR against develop branch`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger = log.New(debug)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

		// Create new branch from base
		currentBranch, err := gitOps.CreateAndSwitchBranch(ctx, baseBranch)
		if err != nil {
			if strings.Contains(err.Error(), "not a git repository") {
				logger.Error("Not in a git repository")
				return fmt.Errorf("please run this command from within a git repository")
			}
			return fmt.Errorf("failed to create branch: %w", err)
		}

		// Always use AI unless title is provided
		autoCommit = prTitle == ""

		// Stage all files first
		if err := gitOps.StageAll(ctx); err != nil {
			if err.Error() == "no changes to commit" {
				logger.Warning("No changes to commit")
				return fmt.Errorf("no changes to create PR from")
			}
			return fmt.Errorf("failed to stage files: %w", err)
		}

		// Generate commit message if in auto mode
		if autoCommit {
			diff, err := gitOps.GetDiff(ctx)
			if err != nil {
				return fmt.Errorf("failed to get diff: %w", err)
			}

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
			commitMsg = "update"
		}

		// Commit changes
		if err := gitOps.Commit(ctx, commitMsg); err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}

		// Push changes
		if err := gitOps.Push(ctx, "origin", currentBranch); err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}

		// Create pull request
		if prTitle == "" {
			prTitle = commitMsg
		}
		if prBody == "" {
			prBody = "Created with ghquick"
		}

		pr, err := ghClient.CreatePullRequest(ctx, repoName, prTitle, prBody, currentBranch, baseBranch)
		if err != nil {
			return err
		}

		logger.Success("🚀 Pull request #%d created successfully!", pr.GetNumber())
		logger.Info("View it here: %s", pr.GetHTMLURL())
		return nil
	},
}

var mergePRCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge a pull request",
	Long: `Merge a pull request by its number or select from a list.
Example:
  ghquick pr merge --number 123
  ghquick pr merge --name my-repo --number 456
  ghquick pr merge  # Interactive selection if multiple PRs exist`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger = log.New(debug)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

		// Initialize GitHub client
		ghClient := github.NewClient(cfg.GitHubToken, debug)

		// If PR number is not provided, list PRs and let user select
		if prNumber == 0 {
			prs, err := ghClient.ListPullRequests(ctx, repoName)
			if err != nil {
				logger.Error("Failed to list pull requests")
				return err
			}

			if len(prs) == 0 {
				logger.Error("No open pull requests found")
				return fmt.Errorf("no open pull requests found")
			}

			if len(prs) == 1 {
				// If only one PR, use it automatically
				prNumber = prs[0].GetNumber()
				logger.Info("Found single PR #%d: %s", prNumber, prs[0].GetTitle())
			} else {
				// Print PRs for selection
				logger.Info("Select a pull request to merge:")
				for i, pr := range prs {
					logger.Info("%d. #%d: %s", i+1, pr.GetNumber(), pr.GetTitle())
				}

				// Read user selection
				var selection int
				fmt.Print("Enter number (1-" + fmt.Sprint(len(prs)) + "): ")
				_, err := fmt.Scanf("%d", &selection)
				if err != nil || selection < 1 || selection > len(prs) {
					logger.Error("Invalid selection")
					return fmt.Errorf("invalid selection")
				}

				prNumber = prs[selection-1].GetNumber()
			}
		}

		// Verify PR exists and is mergeable
		pr, err := ghClient.GetPullRequest(ctx, repoName, prNumber)
		if err != nil {
			logger.Error("Failed to find pull request #%d", prNumber)
			return err
		}

		if !pr.GetMergeable() {
			logger.Error("Pull request #%d cannot be merged", prNumber)
			return fmt.Errorf("pull request is not mergeable")
		}

		// Merge the pull request
		if err := ghClient.MergePullRequest(ctx, repoName, prNumber); err != nil {
			return err
		}

		logger.Success("🎉 Pull request #%d merged successfully!", prNumber)
		return nil
	},
}
