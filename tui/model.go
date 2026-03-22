package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type gradientTickMsg time.Time

type Model struct {
	choices   []string
	cursor    int
	screen    string
	textInput textinput.Model
	result    string
	errorMsg  string
	resultView viewport.Model
	hueShift  float64
	quitting  bool
}

func InitialModel() Model {
	input := textinput.New()
	input.Placeholder = "example: test.txt"
	input.Prompt = ">_< "
	input.CharLimit = 512
	input.Width = 60
	input.Blur()

	vp := viewport.New(80, 20)
	vp.SetContent("")

	return Model{
		choices: []string{
			"Read File (readFile)",
			"Exit",
		},
		screen:     "menu",
		textInput:  input,
		resultView: vp,
	}
}

func (m Model) Init() tea.Cmd {
	return gradientTick()
}
