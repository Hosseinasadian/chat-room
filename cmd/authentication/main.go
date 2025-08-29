package main

import (
	"github.com/hosseinasadian/chat-application/cmd/authentication/command"
	"os"
)

func main() {
	if err := command.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
