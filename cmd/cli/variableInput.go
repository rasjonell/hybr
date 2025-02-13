package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateVariableInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	vars := m.getCurrentVariables()

	cmds := []tea.Cmd{}
	currentVar := vars[m.cursor]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if currentVar.Input.Value() == "" {
				currentVar.Input.SetValue(currentVar.Default)
			}
			if m.cursor == len(vars)-1 {
				if m.activeServiceIndex == len(m.selected)-1 {
					m.step = StepConfirmation
				} else {
					m.activeServiceIndex++
					m.cursor = 0
				}
			} else {
				currentVar.Input.Blur()
				m.cursor++
				cmds = append(cmds, vars[m.cursor].Input.Focus())
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	currentVar.Input, cmd = currentVar.Input.Update(msg)

	return m, tea.Batch(append(cmds, cmd)...)
}

func (m *Model) viewVariableInput() string {
	var lines []string = []string{
		fmt.Sprintf("Please Insert Variables for %s\n", m.currentServiceName()),
	}

	for i, v := range m.getCurrentVariables() {
		lines = append(lines, v.Description, fmt.Sprintf(
			"%d. %s = %s", i+1, v.Name, v.Input.View(),
		))
	}

	return strings.Join(lines, "\n")
}
