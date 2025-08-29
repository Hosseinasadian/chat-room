package main

import (
	AuthenticationCommand "github.com/hosseinasadian/chat-application/cmd/authentication/command"
	UserCommand "github.com/hosseinasadian/chat-application/cmd/user/command"
	"github.com/hosseinasadian/chat-application/pkg/tui"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "chat-app",
	Short: "Main CLI entrypoint",
	Long:  "CLI to manage and run all services (food, taxi, authentication, etc).",
}

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Interactive CI for all services",
}

func main() {
	ciCmd.AddCommand(AuthenticationCommand.RootCommand)
	ciCmd.AddCommand(UserCommand.RootCommand)

	var tuiCmd = &cobra.Command{
		Use:   "tui",
		Short: "Interactive TUI for all services",
		Run: func(cmd *cobra.Command, args []string) {
			tui.RunTUI(ciCmd)
		},
	}

	rootCmd.AddCommand(ciCmd)
	rootCmd.AddCommand(tuiCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
