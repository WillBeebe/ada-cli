package cmd

import (
	"github.com/container-labs/ada/internal/chat"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(chatCmd)
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Get advice from Google's Gemini Ultra LLM",
	Run: func(cmd *cobra.Command, args []string) {
		chat.Chat()
	},
}
