package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateBaseConfigInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	currentVar := m.baseConfigVariables[m.cursor]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, getVariableInputKeys(currentVar.Input.Value()).CtrlA):
			if currentVar.Input.Value() == "" {
				currentVar.Input.SetValue(currentVar.Default)
			}
			for i := m.cursor + 1; i < len(m.baseConfigVariables); i++ {
				m.baseConfigVariables[i].Input.SetValue(m.baseConfigVariables[i].Default)
			}
			m.initServiceSelection()

		case key.Matches(msg, getVariableInputKeys(currentVar.Input.Value()).Return):
			if currentVar.Input.Value() == "" {
				currentVar.Input.SetValue(currentVar.Default)
			}
			if m.cursor == len(m.baseConfigVariables)-1 {
				m.initServiceSelection()
			} else {
				currentVar.Input.Blur()
				m.cursor++
				cmds = append(cmds, m.baseConfigVariables[m.cursor].Input.Focus())
			}

		case key.Matches(msg, getVariableInputKeys(currentVar.Input.Value()).Quit):
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

	currentVarValue := m.baseConfigVariables[m.cursor].Input.Value()
	lines = append(lines, "\n", variableInputHelp.View(getVariableInputKeys(currentVarValue)))

	return strings.Join(lines, "\n")

}
