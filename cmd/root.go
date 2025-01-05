package cmd

import (
	"fmt"

	"github.com/saint0x/ghquick-cli/internal/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghquick",
	Short: "GitHub workflow automation",
	Long: fmt.Sprintf(`
                    %s

%s
  %s                    %s
  %s  %s
  %s          %s
  %s        %s

%s
  %s              %s
  %s      %s
  %s       %s
  %s       %s

  %s              %s
  %s     %s

%s
  %s     %s
  %s    %s
  %s  %s

%s
  1. %s         %s
  2. %s    %s
  3. %s     %s

%s
  %s
  %s
  %s

%s
  %s
  %s
  %s`,
		log.Title("[ ghquick ]"),

		log.Section("PUSH CHANGES"),
		log.Command("ghquick push"), log.Description("AI commit message (default)"),
		log.Command("ghquick push --commitmsg \"fix\""), log.Description("Custom commit message"),
		log.Command("ghquick push --private"), log.Description("Create private repository"),
		log.Command("ghquick push --name repo"), log.Description("Push to specific repo"),

		log.Section("PULL REQUESTS"),
		log.Command("ghquick pr create"), log.Description("AI title & auto-branch (default)"),
		log.Command("ghquick pr create --title"), log.Description("Custom PR title"),
		log.Command("ghquick pr create --body"), log.Description("PR description"),
		log.Command("ghquick pr create --base"), log.Description("Target branch (default: main)"),

		log.Command("ghquick pr merge"), log.Description("Interactive selection (default)"),
		log.Command("ghquick pr merge --number"), log.Description("Merge specific PR by number"),

		log.Section("GLOBAL OPTIONS"),
		log.Command("--name"), log.Description("Target repository (default: current dir)"),
		log.Command("--debug"), log.Description("Show detailed logs"),
		log.Command("--timeout"), log.Description("Operation timeout (default: 2m)"),

		log.Section("QUICK START"),
		log.Command("ghquick push"), log.Description("→ AI generates commit"),
		log.Command("ghquick pr create"), log.Description("→ Creates PR with AI title"),
		log.Command("ghquick pr merge"), log.Description("→ Select PR to merge"),

		log.Section("AI FEATURES"),
		log.Description("• Generates concise commit messages (1-3 words)"),
		log.Description("• Interactive PR selection when multiple PRs exist"),
		log.Description("• AI-powered PR titles from commit messages"),

		log.Section("TIPS"),
		log.Description(fmt.Sprintf("• Use %s for monorepo projects", log.Command("--name"))),
		log.Description(fmt.Sprintf("• Use %s for detailed logs", log.Command("--debug"))),
		log.Description(fmt.Sprintf("• Run %s for shell completion", log.Command("completion [bash|zsh]")))),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	cobra.AddTemplateFunc("section", log.Section)
	cobra.AddTemplateFunc("cmd", log.Command)
	cobra.AddTemplateFunc("desc", log.Description)

	rootCmd.SetUsageTemplate(`{{printf "\n"}}{{section "Usage"}}:
  {{cmd .UseLine}}

{{section "Available Commands"}}:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{cmd .Name | printf "%-11s"}} {{desc .Short}}{{end}}{{end}}

{{section "Flags"}}:{{if .HasAvailableLocalFlags}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces | desc}}{{end}}

{{desc "Use \"ghquick [command] --help\" for more information about a command."}}
`)

	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}
