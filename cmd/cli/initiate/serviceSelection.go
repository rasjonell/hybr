package initiate

import (
	"fmt"
	"github.com/rasjonell/hybr/internal/services"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type selectionKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Space  key.Binding
	Return key.Binding
	Quit   key.Binding
	Help   key.Binding
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
		{k.Up, k.Down, k.Left, k.Right},
		k.ShortHelp(),
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
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "select all"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("␣", "select none"),
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

		case key.Matches(msg, selectionKeys.Left):
			m.selectedServiceNames = []string{}
			m.selected = map[string]services.HybrService{}

		case key.Matches(msg, selectionKeys.Right):
			for _, s := range m.services {
				m.selected[s.GetName()] = s
				m.selectedServiceNames = append(m.selectedServiceNames, s.GetName())
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

			m.cursor = 0
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
	}

	lines = append(lines, "\n", selectionHelp.View(selectionKeys))

	return strings.Join(lines, "\n") + "\n"
}
