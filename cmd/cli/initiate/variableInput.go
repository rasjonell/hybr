package initiate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type variableInputKeyMap struct {
	CtrlA  key.Binding
	Return key.Binding
	Quit   key.Binding
}

var variableInputHelp help.Model = help.New()

func (k variableInputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Return, k.CtrlA, k.Quit}
}

func (k variableInputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func getVariableInputKeys(value string) variableInputKeyMap {
	returnHelp := "save value"
	if value == "" {
		returnHelp = "save default value"
	}

	return variableInputKeyMap{
		CtrlA: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "accept defaults"),
		),
		Return: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", returnHelp),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc/ctrl+c", "quit"),
		),
	}
}

func (m *Model) setDefaultCurrentVar() {
	currentVar := m.getCurrentVariables()[m.cursor]
	if currentVar.Input.Value() == "" {
		currentVar.Input.SetValue(currentVar.Default)
	}
}

func (m *Model) updateVariableInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	vars := m.getCurrentVariables()
	currentVar := vars[m.cursor]

	currentVarValue := currentVar.Input.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, getVariableInputKeys(currentVarValue).CtrlA):
			m.setDefaultCurrentVar()
			for j := m.cursor + 1; j < len(vars); j++ {
				vars[j].Input.SetValue(vars[j].Default)
			}
			m.cursor = 0
			currentVar.Input.Blur()
			if m.activeServiceIndex == len(m.selected)-1 {
				m.step = StepConfirmation
			} else {
				m.activeServiceIndex++
				cmds = append(cmds, m.getCurrentVariables()[m.cursor].Input.Focus())
			}

		case key.Matches(msg, getVariableInputKeys(currentVarValue).Return):
			m.setDefaultCurrentVar()
			if m.cursor == len(vars)-1 {
				if m.activeServiceIndex == len(m.selected)-1 {
					m.step = StepConfirmation
				} else {
					m.cursor = 0
					m.activeServiceIndex++
					currentVar.Input.Blur()
					cmds = append(cmds, m.getCurrentVariables()[m.cursor].Input.Focus())
				}
			} else {
				currentVar.Input.Blur()
				m.cursor++
				cmds = append(cmds, vars[m.cursor].Input.Focus())
			}

		case key.Matches(msg, getVariableInputKeys(currentVarValue).Quit):
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

	currentVarValue := m.getCurrentVariables()[m.cursor].Input.Value()
	lines = append(lines, "\n", variableInputHelp.View(getVariableInputKeys(currentVarValue)))

	return strings.Join(lines, "\n")
}
