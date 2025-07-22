// Package client provides a client for the chat application.
package client

import (
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const ApiUrl = ":4040"

type model struct {
	conn   net.Conn
	loader spinner.Model
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	s.Tick()

	return model{
		loader: s,
	}
}

func connect() tea.Msg {
	conn, err := net.Dial("tcp", ApiUrl)
	if err != nil {
		panic(err)
	}
	return conn
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.loader.Tick, connect)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			return m, nil
		}

	case net.Conn:
		m.conn = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.loader, cmd = m.loader.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.conn == nil {
		return m.loader.View()
	}

	return "Connected"
}

func Start() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
