package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateVariableInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			vars := *m.currentSelectedServiceVariables()
			vars[m.cursor].Value = m.textInput.Value()
			m.textInput.SetValue("")

			if m.cursor == len(vars)-1 {
				if m.activeServiceIndex == len(m.selected)-1 {
					m.step = StepConfirmation
				} else {
					m.activeServiceIndex++
					m.cursor = 0
				}
			} else {
				m.cursor++
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m *Model) viewVariableInput() string {
	var lines []string = []string{
		fmt.Sprintf("Please Insert Variables for %s\n", m.currentServiceName()),
	}

	for i, v := range *m.currentSelectedServiceVariables() {
		textInputView := v.Value
		isSelected := i == m.cursor
		if isSelected {
			textInputView = m.textInput.View()
		}

		lines = append(lines, fmt.Sprintf(
			"%d. %s = %s", i+1, v.Key, textInputView,
		))
	}

	return strings.Join(lines, "\n")
}

func generateVariableInput() textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Focus()

	return ti
}
