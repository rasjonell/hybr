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
			m.selected = map[string]services.SelectedServiceModel{}

		case "right", "l":
			for _, s := range m.services {
				m.selected[s.Name] = services.SelectedServiceModel{Service: &s, Variables: buildVariableList(&s)}
				m.selectedServiceNames = append(m.selectedServiceNames, s.Name)
			}

		case " ":
			selected := m.services[m.cursor]
			if _, exists := m.selected[selected.Name]; exists {
				delete(m.selected, selected.Name)
				m.selectedServiceNames = slices.DeleteFunc(m.selectedServiceNames, func(n string) bool {
					return n == selected.Name
				})
			} else {
				m.selected[selected.Name] = services.SelectedServiceModel{Service: &selected, Variables: buildVariableList(&selected)}
				m.selectedServiceNames = append(m.selectedServiceNames, selected.Name)
			}

		case "enter":
			if len(m.selectedServiceNames) == 0 {
				return m, tea.Quit
			}

			m.cursor = 0
			m.step = StepVariableInput
			m.textInput = generateVariableInput()
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
		if _, exists := m.selected[service.Name]; exists {
			selectionIndicator = "  [x]"
		}

		lines = append(lines, fmt.Sprintf(
			"%s%s %s",
			cursorIndicator, selectionIndicator, service.Description),
		)
	}

	return strings.Join(lines, "\n") + "\n"
}

func buildVariableList(s *services.Service) []services.VariableDefinition {
	varList := make([]services.VariableDefinition, len(s.Variables), len(s.Variables))
	for i, v := range s.Variables {
		varList[i] = services.VariableDefinition{
			Key:   v.Name,
			Value: "",
		}
	}
	return varList
}
