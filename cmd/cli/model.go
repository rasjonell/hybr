package main

import (
	"hybr/internal/services"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Step int

const (
	StepServiceSelection Step = iota
	StepVariableInput
	StepConfirmation
)

type Model struct {
	step Step

	selectedServiceNames []string
	services             []services.Service
	selected             map[string]services.SelectedServiceModel

	cursor             int
	activeServiceIndex int

	textInput textinput.Model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case StepServiceSelection:
		return m.updateServiceSelection(msg)

	case StepVariableInput:
		return m.updateVariableInput(msg)

	case StepConfirmation:
		return m.updateConfirmation(msg)
	}

	return m, nil
}

func (m *Model) View() string {
	switch m.step {
	case StepServiceSelection:
		return m.viewServiceSelection()

	case StepVariableInput:
		return m.viewVariableInput()

	case StepConfirmation:
		return m.viewConfirmation()

	default:
		return "Unknown Step"
	}
}

func (m *Model) currentServiceName() string {
	return m.selectedServiceNames[m.activeServiceIndex]
}

func (m *Model) currentSelectedServiceVariables() *[]services.VariableDefinition {
	v, exists := m.selected[m.selectedServiceNames[m.activeServiceIndex]]
	if !exists {
		panic("fucky wucky")
	}

	return &v.Variables
}

func NewProgram() *tea.Program {
	return tea.NewProgram(&Model{
		step:     StepServiceSelection,
		services: services.GetRegisteredServices(),
		selected: make(map[string]services.SelectedServiceModel),
	})
}
