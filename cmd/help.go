package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var allCommandsCmd = &cobra.Command{
	Use:   "all-commands",
	Short: "Display help for all available commands",
	Long:  `This command shows detailed help information for all commands available in the Ada CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		displayAllCommandsHelp(rootCmd)
	},
}

func init() {
	rootCmd.AddCommand(allCommandsCmd)
}

func displayAllCommandsHelp(cmd *cobra.Command) {
	fmt.Println(titleStyle.Render("Ada CLI - All Commands Reference"))
	fmt.Println(descriptionStyle.Render("This document provides detailed information about all available commands in the Ada CLI."))
	fmt.Println()

	displayCommandHelp(cmd, 0)
}

func displayCommandHelp(cmd *cobra.Command, depth int) {
	if !cmd.IsAvailableCommand() || cmd.IsAdditionalHelpTopicCommand() {
		return
	}

	indent := strings.Repeat("  ", depth)

	fmt.Println(sectionStyle.Render(indent + cmd.Name()))
	fmt.Println(boxStyle.Render(indent + "  " + cmd.Short))

	if cmd.Long != "" {
		fmt.Println(descriptionStyle.Render(indent + "  " + cmd.Long))
	}

	if cmd.UsageString() != "" {
		fmt.Println(infoStyle.Render(indent + "  Usage:"))
		fmt.Println(boxStyle.Render(indent + "    " + cmd.UsageString()))
	}

	if len(cmd.Aliases) > 0 {
		fmt.Println(infoStyle.Render(indent + "  Aliases:"))
		fmt.Println(boxStyle.Render(indent + "    " + strings.Join(cmd.Aliases, ", ")))
	}

	if cmd.HasAvailableFlags() {
		fmt.Println(infoStyle.Render(indent + "  Flags:"))
		buf := new(bytes.Buffer)
		cmd.Flags().PrintDefaults()
		fmt.Println(boxStyle.Render(indent + "    " + replaceFlags(buf.String())))
	}

	if cmd.HasAvailableSubCommands() {
		fmt.Println(infoStyle.Render(indent + "  Subcommands:"))
		for _, subCmd := range cmd.Commands() {
			displayCommandHelp(subCmd, depth+1)
		}
	}

	fmt.Println()
}
