package command

import "github.com/spf13/cobra"

var RootCommand = &cobra.Command{
	Use:   "user",
	Short: "A CLI for user Service",
	Long: `user Service CLI is a tool to manage and run 
the user service, including migrations and server startup.`,
}
