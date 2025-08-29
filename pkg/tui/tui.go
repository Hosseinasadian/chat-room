package tui

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type state int

const (
	selectService state = iota
	selectCommand
)

type model struct {
	root     *cobra.Command
	cancel   context.CancelFunc
	state    state
	services []*cobra.Command
	commands []*cobra.Command
	cursor   int
	selected *cobra.Command
}

// Styles using lipgloss
var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFA500"))
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	highlight    = lipgloss.NewStyle().Background(lipgloss.Color("#333333")).Foreground(lipgloss.Color("#FFFFFF"))
	instruction  = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#888888"))
	selectedText = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")).Bold(true)
)

func NewModel(root *cobra.Command) model {
	return model{
		root:     root,
		state:    selectService,
		services: root.Commands(),
		cursor:   0,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.state == selectCommand {
				// go back to service selection
				m.state = selectService
				m.cursor = 0
				m.selected = nil
				m.commands = nil
			} else {
				// if already at service list, maybe quit?
				return m, tea.Quit
			}
		case "ctrl+c":
			if m.state == selectCommand {
				if m.cancel != nil {
					m.cancel() // graceful shutdown
				}
				return m, tea.Quit
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.state == selectService && m.cursor < len(m.services)-1 {
				m.cursor++
			} else if m.state == selectCommand && m.cursor < len(m.commands)-1 {
				m.cursor++
			}
		case "enter":
			if m.state == selectService {
				m.selected = m.services[m.cursor]
				m.commands = m.selected.Commands()
				m.cursor = 0
				m.state = selectCommand
			} else if m.state == selectCommand {
				cmd := m.commands[m.cursor]
				fmt.Printf("Running: %s %s\n\n", m.selected.Name(), cmd.Name())

				// create context so you can cancel with ctrl+c
				ctx, cancel := context.WithCancel(context.Background())
				m.cancel = cancel
				cmd.SetContext(ctx)

				go func(c *cobra.Command) {
					// If command uses Args, set them here if needed
					c.SetArgs([]string{})

					// If RunE is defined, prefer that
					if c.RunE != nil {
						if err := c.RunE(c, c.Flags().Args()); err != nil {
							fmt.Println("command error:", err)
						}
					} else if c.Run != nil {
						c.Run(c, c.Flags().Args())
					} else {
						fmt.Println("no runnable handler for", c.Name())
					}
				}(cmd)
			}
		case "esc", "left", "b":
			if m.state == selectCommand {
				m.state = selectService
				m.cursor = 0
				m.selected = nil
				m.commands = nil
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "\n"
	if m.state == selectService {
		s += titleStyle.Render("Select a Service") + "\n\n"
		for i, svc := range m.services {
			line := svc.Name()
			if m.cursor == i {
				line = cursorStyle.Render("> " + line)
			} else {
				line = "  " + line
			}
			s += line + "\n"
		}
	} else if m.state == selectCommand {
		s += titleStyle.Render(fmt.Sprintf("Service: %s", m.selected.Name())) + "\n\n"
		s += "Select a Command:\n\n"
		for i, cmd := range m.commands {
			line := cmd.Name()
			if m.cursor == i {
				line = highlight.Render("> " + line)
			} else {
				line = "  " + line
			}
			s += line + "\n"
		}
	}
	s += "\n" + instruction.Render("(Use ↑/↓ to navigate, enter to select, b to go back, q to quit)") + "\n"
	return s
}

func RunTUI(root *cobra.Command) {
	p := tea.NewProgram(NewModel(root), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
