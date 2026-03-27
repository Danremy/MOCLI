package tui

import (
	"os"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type gradientTickMsg time.Time

type Model struct {
	choices    []string
	cursor     int
	screen     string
	filePicker filepicker.Model
	result     string
	errorMsg   string
	resultView viewport.Model
	hueShift   float64
	quitting   bool
}

func InitialModel() Model {
	picker := filepicker.New()
	picker.ShowHidden = false
	picker.ShowPermissions = false
	picker.ShowSize = true
	picker.FileAllowed = true
	picker.DirAllowed = false
	picker.CurrentDirectory = "."
	picker.SetHeight(12)

	if cwd, err := os.Getwd(); err == nil {
		picker.CurrentDirectory = cwd
	}

	vp := viewport.New(80, 20)
	vp.SetContent("")

	return Model{
		choices: []string{
			"Read File (readFile)",
			"Exit",
		},
		screen:     "menu",
		filePicker: picker,
		resultView: vp,
	}
}

func (m Model) Init() tea.Cmd {
	// Initialize picker state early; it will also be refreshed when entering picker screen.
	return tea.Batch(gradientTick(), m.filePicker.Init())
}
