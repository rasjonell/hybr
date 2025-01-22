package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateConfirmation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) viewConfirmation() string {
	return fmt.Sprintf(
		"Press Enter to install %d service(s): %s",
		len(m.selectedServiceNames), strings.Join(m.selectedServiceNames, ", "),
	)
}
