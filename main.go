package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Message struct {
	msg string
}

type Model struct {
	messages []Message
	promt    string
}

func initalModel() *Model {
	m := &Model{
		messages: make([]Message, 0),
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.messages = append(m.messages, Message{msg: strings.TrimSpace(m.promt) + "\n"})
			m.promt = ""
		case "backspace":
			if len(m.promt) > 0 {
				m.promt = m.promt[:len(m.promt)-1]
			}
		default:
			m.promt += string(msg.String())
		}
	}

	return m, nil
}

func (m Model) View() string {
	heading := lipgloss.NewStyle().Bold(true).PaddingBottom(1).Render("Welcome to the chat!")

	messages := make([]string, 0)
	for _, message := range m.messages {
		messages = append(messages, message.msg)
	}

	seperator := lipgloss.NewStyle().Render("------------------------------------")

	prompt := lipgloss.NewStyle().Render(m.promt)

	return lipgloss.JoinVertical(lipgloss.Top, heading, lipgloss.JoinVertical(lipgloss.Top, messages...), seperator, prompt)
}

func main() {
	// tcpServer, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6969")
	// if err != nil {
	// 	fmt.Println("Resolve add failed")
	// 	os.Exit(1)
	// }

	// conn, err := net.DialTCP("tcp", nil, tcpServer)
	// if err != nil {
	// 	fmt.Println("Could not connect to server")
	// 	os.Exit(1)
	// }

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(initalModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
