package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#9B9B9B"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Background(lipgloss.Color("#1A1A1A")).
			Bold(true).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1).
			Width(50).
			Align(lipgloss.Center)

	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			UnderlineSpaces(true).
			Underline(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0"))

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4EC9B0")).
			Bold(true)

	flagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500"))

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC")).
				MarginTop(1).
				MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1).
			MarginTop(1).
			MarginBottom(1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Center)
)

var rootCmd = &cobra.Command{
	Use:   "ada",
	Short: "Ada is a CLI tool",
	Long:  "Ada is a CLI tool for various tasks",
}

func init() {
	cobra.AddTemplateFunc("styleTitle", func(s string) string {
		return titleStyle.Render(s)
	})

	cobra.AddTemplateFunc("styleSection", func(s string) string {
		return sectionStyle.Render(s)
	})

	cobra.AddTemplateFunc("styleDescription", func(s string) string {
		return descriptionStyle.Render(s)
	})

	cobra.AddTemplateFunc("styleBox", func(s string) string {
		return boxStyle.Render(s)
	})

	cobra.AddTemplateFunc("styleFooter", func(s string) string {
		return footerStyle.Render(s)
	})

	cobra.AddTemplateFunc("styleCmdName", func(s string) string {
		return commandStyle.Render(s)
	})

	cobra.AddTemplateFunc("styleInfo", func(s string) string {
		return infoStyle.Render(s)
	})

	cobra.AddTemplateFunc("replaceFlags", func(s string) string {
		return replaceFlags(s)
	})

	helpTemplate := `
{{styleTitle "Ada CLI"}}

{{styleDescription .Long}}

{{styleSection "Usage:"}}
{{styleBox (printf "  %s" .UseLine)}}

{{if .HasAvailableSubCommands}}
{{styleSection "Available Commands:"}}
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{styleCmdName .Name | printf "%-15s"}} {{styleInfo .Short}}{{end}}{{end}}
{{end}}

{{if .HasAvailableLocalFlags}}
{{styleSection "Flags:"}}
{{styleBox (.LocalFlags.FlagUsages | trimTrailingWhitespaces | replaceFlags)}}
{{end}}

{{if .HasAvailableInheritedFlags}}
{{styleSection "Global Flags:"}}
{{styleBox (.InheritedFlags.FlagUsages | trimTrailingWhitespaces | replaceFlags)}}
{{end}}

{{if .HasExample}}
{{styleSection "Examples:"}}
{{styleBox .Example}}
{{end}}

{{styleFooter "For more information, visit: https://github.com/your-repo/ada"}}
`

	rootCmd.SetHelpTemplate(helpTemplate)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func replaceFlags(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if strings.Contains(line, "--") {
			parts := strings.SplitN(line, "   ", 2)
			if len(parts) == 2 {
				lines[i] = flagStyle.Render(parts[0]) + "   " + infoStyle.Render(parts[1])
			}
		}
	}
	return strings.Join(lines, "\n")
}
