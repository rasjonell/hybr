package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type confirmationKeyMap struct {
	Return key.Binding
	Quit   key.Binding
}

var confirmationHelp help.Model = help.New()

func (k confirmationKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Return, k.Quit}
}

func (k confirmationKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var confirmationKeys = confirmationKeyMap{
	Return: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "confirm installation"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

func (m *Model) updateConfirmation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, selectionKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, confirmationKeys.Return):
			m.buildFinalServices()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) viewConfirmation() string {
	lines := []string{
		fmt.Sprintf(
			"Press Enter to install %d service(s): %s",
			len(m.selectedServiceNames), strings.Join(m.selectedServiceNames, ", "),
		),
		"\n",
		confirmationHelp.View(confirmationKeys),
	}

	return strings.Join(lines, "\n") + "\n"
}
