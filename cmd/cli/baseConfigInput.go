package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateBaseConfigInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	currentVar := m.baseConfigVariables[m.cursor]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if currentVar.Input.Value() == "" {
				currentVar.Input.SetValue(currentVar.Default)
			}
			if m.cursor == len(m.baseConfigVariables)-1 {
				m.cursor = 0
				m.step = StepServiceSelection
			} else {
				currentVar.Input.Blur()
				m.cursor++
				cmds = append(cmds, m.baseConfigVariables[m.cursor].Input.Focus())
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	currentVar.Input, cmd = currentVar.Input.Update(msg)

	return m, tea.Batch(append(cmds, cmd)...)
}

func (m *Model) viewBaseConfigInput() string {
	var lines []string

	for i, v := range m.baseConfigVariables {
		lines = append(lines, v.Description, fmt.Sprintf(
			"%d. %s = %s\n", i+1, v.Name, v.Input.View(),
		))
	}

	return strings.Join(lines, "\n")

}
