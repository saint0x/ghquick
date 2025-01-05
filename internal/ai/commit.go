package ai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type GenerateResult struct {
	Message string
	Error   error
}

type CommitMessageGenerator struct {
	client *openai.Client
}

func NewCommitMessageGenerator(apiKey string) *CommitMessageGenerator {
	return &CommitMessageGenerator{
		client: openai.NewClient(apiKey),
	}
}

const systemPrompt = `You are a commit message generator. Your task is to create concise, descriptive commit messages that are 1-3 words long.
Rules:
1. Be extremely concise and straightforward
2. Focus on the core change or feature
3. No fluff or unnecessary words
4. Dark humor or inappropriate jokes are allowed but not required
5. No emojis or special characters
6. No conventional commit prefixes (feat:, fix:, etc.)

Examples:
- "fix memory leak"
- "add user auth"
- "optimize queries"
- "nuke legacy code"
- "unfuck database"
- "add tests"`

func (g *CommitMessageGenerator) GenerateFromDiffAsync(ctx context.Context, diff string, result chan<- GenerateResult) {
	go func() {
		resp, err := g.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: "gpt-4o-mini",
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
				MaxTokens:   10,
				Temperature: 0.7,
			},
		)

		if err != nil {
			result <- GenerateResult{Error: fmt.Errorf("failed to generate commit message: %w", err)}
			return
		}

		if len(resp.Choices) == 0 {
			result <- GenerateResult{Error: fmt.Errorf("no commit message generated")}
			return
		}

		result <- GenerateResult{Message: resp.Choices[0].Message.Content}
	}()
}
