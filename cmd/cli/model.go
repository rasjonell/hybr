package main

import (
	"hybr/internal/nginx"
	"hybr/internal/services"

	tea "github.com/charmbracelet/bubbletea"
)

type Step int

const (
	StepBaseConfigInput Step = iota
	StepServiceSelection
	StepVariableInput
	StepConfirmation
)

type Model struct {
	step Step

	selectedServiceNames []string
	services             []services.HybrService
	selected             map[string]services.HybrService

	finalBaseConfig     nginx.NginxConfig
	baseConfigVariables []*services.VariableDefinition

	cursor             int
	activeServiceIndex int
}

var model *Model

func InitCLI() {
	registeredServices := services.GetRegisteredServices()

	step := StepBaseConfigInput
	if flags.isBaseConfigComplete || flags.forceNoSSL {
		step = StepServiceSelection
	}

	model = &Model{
		cursor:              0,
		step:                step,
		services:            registeredServices,
		baseConfigVariables: getBaseConfigVariables(),
		selected:            make(map[string]services.HybrService),
	}
}

func (m *Model) initInputs() {
	focusTaken := false
	for _, serviceName := range m.selectedServiceNames {
		s := m.selected[serviceName]
		for template, variableDefinitions := range s.GetVariables() {
			for _, v := range variableDefinitions {
				ti := buildTextInput(v.Default)
				if !focusTaken {
					ti.Focus()
					focusTaken = true
				}
				v.Input = ti
				v.Value = ""
				v.Template = template
			}
		}
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	default:
		switch m.step {
		case StepBaseConfigInput:
			return m.updateBaseConfigInput(msg)

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
	case StepBaseConfigInput:
		return m.viewBaseConfigInput()

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

func (m *Model) getCurrentSelectedService() services.HybrService {
	return m.selected[m.selectedServiceNames[m.activeServiceIndex]]
}

func (m *Model) getCurrentVariables() []*services.VariableDefinition {
	defs := make([]*services.VariableDefinition, 0)
	for _, vars := range m.getCurrentSelectedService().GetVariables() {
		defs = append(defs, vars...)
	}
	return defs
}

func (m *Model) buildFinalServices() {
	for _, service := range m.selected {
		for _, vars := range service.GetVariables() {
			for _, v := range vars {
				v.Value = v.Input.Value()
			}
		}
	}

	if flags.isBaseConfigComplete {
		m.finalBaseConfig = &nginx.ConfigImpl{
			Email:  flags.email,
			Domain: flags.domain,
		}
		return
	}

	var finalBaseConfig nginx.ConfigImpl
	for _, def := range m.baseConfigVariables {
		switch def.Name {
		case "Email":
			val := def.Input.Value()
			if val == "" {
				val = def.Default
			}
			finalBaseConfig.Email = val
		case "Domain":
			val := def.Input.Value()
			if val == "" {
				val = def.Default
			}
			finalBaseConfig.Domain = val
		}
	}
	m.finalBaseConfig = &finalBaseConfig
}

func (m *Model) getFinalServices() []services.HybrService {
	finalServices := make([]services.HybrService, 0)
	for _, s := range m.selected {
		finalServices = append(finalServices, s)
	}
	return finalServices
}

func NewProgram() *tea.Program {
	return tea.NewProgram(model)
}
