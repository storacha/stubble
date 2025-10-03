package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	text string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyRunes {
			m.text += msg.String()
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "Hello, World!!!\n\n> " + m.text
}
