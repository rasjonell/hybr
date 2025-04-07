package initiate

import (
	"fmt"
	"github.com/rasjonell/hybr/internal/services"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type selectionKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Space  key.Binding
	Return key.Binding
	Quit   key.Binding
	Help   key.Binding
	CtrlA  key.Binding
}

var selectionHelp help.Model = help.New()

func (m *Model) initServiceSelection() {
	m.cursor = 0
	m.step = StepServiceSelection
	selectionHelp.ShowAll = true
}

func (k selectionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Space, k.Return, k.Help, k.Quit}
}

func (k selectionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Return, k.Help},
		{k.Space, k.CtrlA, k.Quit},
	}
}

var selectionKeys = selectionKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	CtrlA: key.NewBinding(
		key.WithKeys("ctrl+a", "a"),
		key.WithHelp("ctrl+a", "select all"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("␣", "select"),
	),
	Return: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "finish selection"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func (m *Model) updateServiceSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, selectionKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, selectionKeys.Help):
			selectionHelp.ShowAll = !selectionHelp.ShowAll

		case key.Matches(msg, selectionKeys.Up):
			if m.cursor == 0 {
				m.cursor = len(m.services) - 1
			} else {
				m.cursor--
			}

		case key.Matches(msg, selectionKeys.Down):
			if m.cursor == len(m.services)-1 {
				m.cursor = 0
			} else {
				m.cursor++
			}

		case key.Matches(msg, selectionKeys.CtrlA):
			if len(m.services) == len(m.selectedServiceNames) {
				m.selected = make(map[string]services.HybrService)
				m.selectedServiceNames = make([]string, 0)
			} else {
				for _, s := range m.services {
					m.selected[s.GetName()] = s
					m.selectedServiceNames = append(m.selectedServiceNames, s.GetName())
				}
			}

		case key.Matches(msg, selectionKeys.Space):
			selected := m.services[m.cursor]
			if _, exists := m.selected[selected.GetName()]; exists {
				delete(m.selected, selected.GetName())
				m.selectedServiceNames = slices.DeleteFunc(m.selectedServiceNames, func(n string) bool {
					return n == selected.GetName()
				})
			} else {
				m.selected[selected.GetName()] = selected
				m.selectedServiceNames = append(m.selectedServiceNames, selected.GetName())
			}

		case key.Matches(msg, selectionKeys.Return):
			if len(m.selectedServiceNames) == 0 {
				return m, tea.Quit
			}
			m.initInputs()
			m.step = StepVariableInput
		}
	}

	return m, nil
}

func (m *Model) viewServiceSelection() string {
	selectionHelp.Width = m.weight
	var lines []string = []string{
		"Please Select The Services You Want To Install:\n",
	}

	for i, service := range m.services {
		_, exists := services.GetRegistry().GetInstallation(service.GetName())

		cursorIndicator := " "
		if i == m.cursor {
			cursorIndicator = ">"
		}

		selectionIndicator := "  [ ]"
		if _, exists := m.selected[service.GetName()]; exists {
			selectionIndicator = "  [x]"
		}

		lines = append(lines,
			fmt.Sprintf(
				"%s%s %s",
				cursorIndicator, selectionIndicator, service.GetDescription(),
			),
		)

		if exists {
			style := lipgloss.NewStyle().
				Bold(true).
				Background(lipgloss.Color("#ff0000")).
				MarginLeft(len(selectionIndicator) + 2)

			lines = append(lines,
				style.Render("This is an already installed service, selecting it will override current installation!"),
				"",
			)
		} else {
			style := lipgloss.NewStyle().
				Bold(true).
				Background(lipgloss.Color("#083808")).
				MarginLeft(len(selectionIndicator) + 2)

			lines = append(lines,
				style.Render("This is new service, selecting it will create a fresh new installation!"),
				"",
			)
		}

	}

	lines = append(lines, "\n", selectionHelp.View(selectionKeys))

	return strings.Join(lines, "\n") + "\n"
}
