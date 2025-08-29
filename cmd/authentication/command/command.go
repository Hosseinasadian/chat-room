package command

import "github.com/spf13/cobra"

var RootCommand = &cobra.Command{
	Use:   "authentication",
	Short: "A CLI for authentication Service",
	Long: `authentication Service CLI is a tool to manage and run 
the authentication service, including migrations and server startup.`,
}
