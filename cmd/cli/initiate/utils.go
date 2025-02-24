package initiate

import "github.com/charmbracelet/bubbles/textinput"

func buildTextInput(def string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = def

	return ti
}
