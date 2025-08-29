package main

import (
	"github.com/hosseinasadian/chat-application/cmd/user/command"
	"os"
)

func main() {
	if err := command.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
