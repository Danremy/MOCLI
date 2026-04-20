package app

import (
	"fmt"
	"os"

	"MOCLI/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type readerModel struct {
	explorer components.FileExplorer
	selected string
	quitting bool
}

func (m readerModel) Init() tea.Cmd { return nil }

func (m readerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "backspace", "h":
			m.explorer.GoUp()
			return m, nil
		case "enter":
			if path, ok := m.explorer.Selected(); ok {
				m.selected = path
				return m, tea.Quit
			}
			return m, nil
		}
	}
	return m, m.explorer.Update(msg)
}

func (m readerModel) View() string {
	if m.quitting {
		return ""
	}
	return m.explorer.View() + "\n\nenter 选择 · h 返回上级 · q 退出\n"
}

func SelectFile(startPath string) (string, error) {
	m := readerModel{explorer: components.NewFileExplorer(startPath)}
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	return final.(readerModel).selected, nil
}

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取失败: %w", err)
	}
	return string(content), nil
}
