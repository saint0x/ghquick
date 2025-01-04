package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type CommitMessageGenerator struct {
	client *openai.Client
}

func NewCommitMessageGenerator(apiKey string) *CommitMessageGenerator {
	return &CommitMessageGenerator{
		client: openai.NewClient(apiKey),
	}
}

type GenerateResult struct {
	Message string
	Error   error
}

func (g *CommitMessageGenerator) GenerateFromDiffAsync(ctx context.Context, diff string, resultChan chan<- GenerateResult) {
	go func() {
		message, err := g.GenerateFromDiff(ctx, diff)
		resultChan <- GenerateResult{
			Message: message,
			Error:   err,
		}
	}()
}

func (g *CommitMessageGenerator) GenerateFromDiff(ctx context.Context, diff string) (string, error) {
	systemPrompt := `You are a commit message generator. Given a git diff, generate a concise, 
descriptive commit message following conventional commits format. Focus on the main changes and their purpose.
Format: <type>(<scope>): <description>
Types: feat, fix, docs, style, refactor, test, chore
Keep it under 72 characters.`

	resp, err := g.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4-1106-preview",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Generate a commit message for this diff:\n\n%s", diff),
				},
			},
			MaxTokens:   60,
			Temperature: 0.3,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	message := strings.TrimSpace(resp.Choices[0].Message.Content)
	return message, nil
}
