package main

import (
	"fmt"
	"hybr/internal/services"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateServiceSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor == 0 {
				m.cursor = len(m.services) - 1
			} else {
				m.cursor--
			}

		case "down", "j":
			if m.cursor == len(m.services)-1 {
				m.cursor = 0
			} else {
				m.cursor++
			}

		case "left", "h":
			m.selectedServiceNames = []string{}
			m.selected = map[string]services.HybrService{}

		case "right", "l":
			for _, s := range m.services {
				m.selected[s.GetName()] = s
				m.selectedServiceNames = append(m.selectedServiceNames, s.GetName())
			}

		case " ":
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

		case "enter":
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

		lines = append(lines, fmt.Sprintf(
			"%s%s %s",
			cursorIndicator, selectionIndicator, service.GetDescription()),
		)
	}

	return strings.Join(lines, "\n") + "\n"
}
