package tui

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
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
		if m.screen != "menu" {
			return m, nil
		}
		return updateGradient(m)
	case tea.WindowSizeMsg:
		contentWidth := msg.Width - 4
		if contentWidth < 20 {
			contentWidth = 20
		}
		contentHeight := msg.Height - 8
		if contentHeight < 6 {
			contentHeight = 6
		}
		m.resultView.Width = contentWidth
		m.resultView.Height = contentHeight
		pickerHeight := msg.Height - 12
		if pickerHeight < 8 {
			pickerHeight = 8
		}
		m.filePicker.SetHeight(pickerHeight)
		if m.screen == "picker" {
			var cmd tea.Cmd
			m.filePicker, cmd = m.filePicker.Update(msg)
			return m, cmd
		}
		return m, nil
	case tea.KeyMsg:
		if m.screen == "picker" {
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "esc":
				m.screen = "menu"
				m.errorMsg = ""
				m.result = ""
				m.resultView.SetContent("")
				return m, gradientTick()
			default:
				var cmd tea.Cmd
				m.filePicker, cmd = m.filePicker.Update(msg)

				if ok, path := m.filePicker.DidSelectFile(msg); ok {
					m.readSelectedFile(path)
					m.screen = "result"
					return m, cmd
				}

				if ok, path := m.filePicker.DidSelectDisabledFile(msg); ok {
					m.errorMsg = fmt.Sprintf("该文件类型暂不可读: %s", path)
				}

				return m, cmd
			}
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
				m.resultView.SetContent("")
				return m, gradientTick()
			}
			return m, nil
		case "up", "k":
			if m.screen == "result" {
				var cmd tea.Cmd
				m.resultView, cmd = m.resultView.Update(msg)
				return m, cmd
			}
			if m.screen == "menu" && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.screen == "result" {
				var cmd tea.Cmd
				m.resultView, cmd = m.resultView.Update(msg)
				return m, cmd
			}
			if m.screen == "menu" && m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.screen == "result" {
				return m, nil
			}
			if m.screen != "menu" {
				return m, nil
			}

			m.errorMsg = ""
			m.result = ""

			switch m.cursor {
			case 0:
				m.errorMsg = ""
				if cwd, err := os.Getwd(); err == nil {
					m.filePicker.CurrentDirectory = cwd
				}
				m.screen = "picker"
				return m, m.filePicker.Init()
			case 1:
				m.quitting = true
				return m, tea.Quit
			}
		case "pgup", "pgdown", "u", "d":
			if m.screen == "result" {
				var cmd tea.Cmd
				m.resultView, cmd = m.resultView.Update(msg)
				return m, cmd
			}
		default:
			if m.screen == "result" {
				var cmd tea.Cmd
				m.resultView, cmd = m.resultView.Update(msg)
				return m, cmd
			}
		}
	}

	if m.screen == "picker" {
		// Forward non-key messages (like filepicker readDirMsg) so file list can load.
		var cmd tea.Cmd
		m.filePicker, cmd = m.filePicker.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) readSelectedFile(path string) {
	cleanPath := strings.TrimSpace(path)
	if cleanPath == "" {
		m.errorMsg = "请选择一个文件"
		m.result = ""
		m.resultView.SetContent(m.errorMsg)
		return
	}

	absolutePath, err := filepath.Abs(cleanPath)
	if err != nil {
		absolutePath = cleanPath
	}

	content, err := fileio.ReadFileContent(absolutePath)
	if err != nil {
		m.errorMsg = fmt.Sprintf("读取文件失败: %v", err)
		m.result = ""
		m.resultView.SetContent(m.errorMsg)
		return
	}

	displayContent := strings.ReplaceAll(content, "\r\n", "\n")
	displayContent = strings.ReplaceAll(displayContent, "\r", "\n")
	if strings.TrimSpace(content) == "" {
		displayContent = "（文件为空）"
	}

	m.errorMsg = ""
	m.result = fmt.Sprintf(
		"文件路径: %s\n内容长度: %d bytes\n\n--- 文件内容开始 ---\n%s\n--- 文件内容结束 ---",
		absolutePath,
		len(content),
		displayContent,
	)
	m.resultView.SetContent(m.result)
	m.resultView.GotoTop()
}

func Run() {
	// Use alternate screen for a clean terminal after quitting the TUI.
	program := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "运行 TUI 失败: %v\n", err)
		os.Exit(1)
	}
}
