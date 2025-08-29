package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the user service",
	Long:  `This command starts the user service.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	fmt.Println("Starting user service")
}

func init() {
	RootCommand.AddCommand(serveCmd)
}
