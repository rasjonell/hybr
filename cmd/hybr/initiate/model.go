package initiate

import (
	"github.com/rasjonell/hybr/internal/services"

	tea "github.com/charmbracelet/bubbletea"
)

type Step int

const (
	StepServiceSelection Step = iota
	StepVariableInput
	StepConfirmation
)

type Model struct {
	Done bool

	step Step

	selectedServiceNames []string
	services             []services.HybrService
	selected             map[string]services.HybrService

	cursor             int
	activeServiceIndex int

	height int
	weight int
}

var model *Model

func GetModel() *Model {
	return model
}

func InitCLI() {
	registeredServices := services.GetRegisteredServices()

	model = &Model{
		cursor:   0,
		Done:     false,
		services: registeredServices,
		step:     StepServiceSelection,
		selected: make(map[string]services.HybrService),
	}

	model.initServiceSelection()
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
	return tea.SetWindowTitle("Hybr CLI")
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.weight = msg.Width
		m.height = msg.Height
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
	m.Done = true
	for _, service := range m.selected {
		for _, vars := range service.GetVariables() {
			for _, v := range vars {
				v.Value = v.Input.Value()
			}
		}
	}
}

func (m *Model) GetFinalServices() []services.HybrService {
	finalServices := make([]services.HybrService, 0)
	for _, s := range m.selected {
		finalServices = append(finalServices, s)
	}
	return finalServices
}

func NewProgram() *tea.Program {
	return tea.NewProgram(model, tea.WithAltScreen())
}
