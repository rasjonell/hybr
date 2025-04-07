package initiate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/rasjonell/hybr/internal/services"
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
	newServices := []string{}
	existingServices := []string{}

	for _, name := range m.selectedServiceNames {
		if _, exists := services.GetRegistry().GetInstallation(name); exists {
			existingServices = append(existingServices, name)
		} else {
			newServices = append(newServices, name)
		}
	}

	lines := []string{}

	if len(newServices) != 0 {
		lines = append(lines,
			lipgloss.
				NewStyle().
				Bold(true).
				Background(lipgloss.Color("#083808")).
				Render(fmt.Sprintf("%d New Service(s) Will Be Installed", len(newServices))),
		)

		l := list.New(newServices)

		lines = append(lines, l.String(), "\n")
	}

	if len(existingServices) != 0 {
		lines = append(lines,
			lipgloss.
				NewStyle().
				Bold(true).
				Background(lipgloss.Color("#ff0000")).
				Render(fmt.Sprintf("%d Existing Service(s) Will Be Overridden", len(existingServices))),
		)

		l := list.New(existingServices)

		lines = append(lines, l.String(), "\n")
	}

	lines = append(lines,
		"Press Enter To Confirm\n",
		confirmationHelp.View(confirmationKeys),
	)

	return strings.Join(lines, "\n") + "\n"
}
