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
	Template    string
	Name        string
	Default     string
	Description string
	Input       textinput.Model
}

type ServiceModel struct {
	Name        string
	Description string
	Templates   []string
	Variables   []*Variable
}

type Model struct {
	step Step

	selectedServiceNames []string
	services             []*ServiceModel
	selected             map[string]*ServiceModel
	finalServices        []*services.SelectedServiceModel

	cursor             int
	activeServiceIndex int
}

var model *Model

func init() {
	services.InitRegistry(flags.forceResetTemplates)
	registeredServices := services.GetRegisteredServices()
	modelServices := make([]*ServiceModel, len(registeredServices), cap(registeredServices))

	for i, s := range registeredServices {
		templateCount := 0
		vars := []*Variable{}

		for template, variableDefinitions := range s.Variables {
			for i, v := range variableDefinitions {
				ti := textinput.New()
				ti.Prompt = ""
				ti.Placeholder = v.Default
				if i == templateCount && i == 0 {
					ti.Focus()
				}

				vars = append(vars, &Variable{
					Input:       ti,
					Name:        v.Name,
					Template:    template,
					Default:     v.Default,
					Description: v.Description,
				})
			}
			templateCount++
		}

		modelServices[i] = &ServiceModel{
			Variables:   vars,
			Name:        s.Name,
			Templates:   s.Templates,
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

func (m *Model) buildFinalServices() {
	finalServices := []*services.SelectedServiceModel{}

	for serviceName, service := range m.selected {
		variableDefinitions := make(map[string][]*services.VariableDefinition)

		for _, v := range service.Variables {
			defSlice, exists := variableDefinitions[v.Template]
			varDef := &services.VariableDefinition{
				Key:   v.Name,
				Value: v.Input.Value(),
			}
			if exists {
				variableDefinitions[v.Template] = append(defSlice, varDef)
			} else {
				defs := []*services.VariableDefinition{}
				variableDefinitions[v.Template] = append(defs, varDef)
			}
		}

		finalServices = append(finalServices, &services.SelectedServiceModel{
			ServiceName: serviceName,
			Variables:   variableDefinitions,
		})
	}

	m.finalServices = finalServices
}

func NewProgram() *tea.Program {
	return tea.NewProgram(model)
}
