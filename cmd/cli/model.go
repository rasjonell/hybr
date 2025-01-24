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

type final uint8

type Variable struct {
	Name        string
	Default     string
	Description string
	Input       textinput.Model
}

type ServiceModel struct {
	Name        string
	Description string
	Variables   []*Variable
}

type Model struct {
	step Step

	selectedServiceNames []string
	services             []ServiceModel
	selected             map[string]*ServiceModel

	cursor             int
	activeServiceIndex int
}

var model *Model

func init() {
	services.InitRegistry(flags.forceResetTemplates)
	services := services.GetRegisteredServices()
	modelServices := make([]ServiceModel, len(services), cap(services))

	for i, s := range services {
		vars := make([]*Variable, len(s.Variables), cap(s.Variables))
		for i, v := range s.Variables {
			ti := textinput.New()
			ti.Prompt = ""
			ti.Placeholder = v.Default
			if i == 0 {
				ti.Focus()
			}

			vars[i] = &Variable{
				Input:       ti,
				Name:        v.Name,
				Default:     v.Default,
				Description: v.Description,
			}
		}

		modelServices[i] = ServiceModel{
			Variables:   vars,
			Name:        s.Name,
			Description: s.Description,
		}
	}

	model = &Model{
		services: modelServices,
		step:     StepServiceSelection,
		selected: make(map[string]*ServiceModel),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	default:
		switch m.step {
		case StepServiceSelection:
			return m.updateServiceSelection(msg)

		case StepVariableInput:
			return m.updateVariableInput(msg)

		case StepConfirmation:
			return m.updateConfirmation(msg)
		}
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

func (m *Model) getCurrentSelectedService() *ServiceModel {
	return m.selected[m.selectedServiceNames[m.activeServiceIndex]]
}

func NewProgram() *tea.Program {
	return tea.NewProgram(model)
}
