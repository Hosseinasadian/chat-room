Chat Application CLI

This project provides a unified command-line interface (CLI) for managing and running all services in the chat application (authentication, user, etc). It is built with:

Cobra
→ for CLI commands and flags

Bubbletea
→ for a text-based interactive UI (TUI)

Lipgloss
→ for styled terminal output

The CLI exposes two main entrypoints:

ci → Run commands directly (automation-friendly)

tui → Run services interactively via a text user interface

📦 Running the CLI
1. Run the TUI (interactive mode)

The TUI is a full-screen terminal app built with Bubbletea. It lets you browse available services and commands using the keyboard.

go run cmd/main.go tui


Use ↑ / ↓ to navigate

Press Enter to select a service or run a command

Press b or Esc to go back

Press q to quit

Example:

Run go run cmd/main.go tui

Select authentication service

Pick the serve command → this starts the Authentication service with live output inside your terminal

2. Run commands directly (CI mode)

You can also run services directly from the CLI without going through the TUI. This is especially useful for automation, scripts, and CI/CD pipelines.

go run cmd/main.go ci <service> <command> [flags]


Example: run the authentication service:

go run cmd/main.go ci authentication serve


This directly executes the same command that the TUI would trigger, but without the interactive interface.

⚙️ How it works

cmd/main.go registers all services (authentication, user, etc.) as subcommands of ci

tui starts a Bubbletea program that shows a menu of available services and commands

Selecting a command in the TUI runs the same Cobra command handler as if it were called via ci

So both tui and ci ultimately run the same underlying Cobra commands.

🧑‍💻 Development Notes

Add new services under cmd/<service>/command

Each service should expose a RootCommand (a Cobra command) with its own subcommands (like serve, migrate, etc.)

The TUI automatically discovers services and their commands from the Cobra tree

Configuration is loaded via configloader using YAML + environment variables

🖥 Example Session

Interactive (TUI):

$ go run cmd/main.go tui

Select a Service
> authentication
user

# Enter → shows available commands
Service: authentication
Select a Command:

> serve
migrate


Direct CLI (CI):

$ go run cmd/main.go ci authentication serve
Running authentication service...


🔑 In short:

Use tui if you want an interactive dashboard-like interface

Use ci if you want scripting, automation, or CI/CD compatibility