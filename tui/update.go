package tui

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"MOCLI/internal/fileio"

	tea "github.com/charmbracelet/bubbletea"
)

// #region  渐变背景相关代码

func gradientTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return gradientTickMsg(t)
	})
}
func updateGradient(m Model) (Model, tea.Cmd) {
	m.hueShift = math.Mod(m.hueShift+0.18, math.Pi*2)
	return m, gradientTick()
}

// #endregion

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gradientTickMsg:
		return updateGradient(m)
	case tea.KeyMsg:
		if m.screen == "input" {
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "esc":
				m.screen = "menu"
				m.errorMsg = ""
				m.result = ""
			case "enter":
				path := strings.TrimSpace(m.textInput.Value())
				if path == "" {
					m.errorMsg = "请输入文件路径，例如 test.txt"
					m.result = ""
					m.screen = "result"
					return m, nil
				}

				content, err := fileio.ReadFileContent(path)
				if err != nil {
					m.errorMsg = fmt.Sprintf("读取文件失败: %v", err)
					m.result = ""
				} else {
					displayContent := content
					if strings.TrimSpace(content) == "" {
						displayContent = "（文件为空）"
					}
					m.errorMsg = ""
					m.result = fmt.Sprintf(
						"文件路径: %s\n内容长度: %d bytes\n\n--- 文件内容开始 ---\n%s\n--- 文件内容结束 ---",
						path,
						len(content),
						displayContent,
					)
				}
				m.screen = "result"
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}

			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.screen == "result" {
				m.screen = "menu"
				m.result = ""
				m.errorMsg = ""
			}
		case "up", "k":
			if m.screen == "menu" && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.screen == "menu" && m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.screen != "menu" {
				return m, nil
			}

			m.errorMsg = ""
			m.result = ""

			switch m.cursor {
			case 0:
				m.textInput.SetValue("")
				m.textInput.Focus()
				m.screen = "input"
			case 1:
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func Run() {
	program := tea.NewProgram(InitialModel())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "运行 TUI 失败: %v\n", err)
		os.Exit(1)
	}
}
