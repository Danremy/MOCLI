package cmd

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type lessonModel struct {
	items    []string
	cursor   int
	screen   string
	input    string
	result   string
	errorMsg string
	hueShift float64
	quitting bool
}

type rainbowTickMsg time.Time

func initialLessonModel() lessonModel {
	return lessonModel{
		items: []string{
			"读取文件（readFile）",
			"退出",
		},
		screen: "menu",
		input:  "test.txt",
	}
}

func (m lessonModel) Init() tea.Cmd {
	return rainbowTick()
}

func rainbowTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return rainbowTickMsg(t)
	})
}

func (m lessonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case rainbowTickMsg:
		m.hueShift = math.Mod(m.hueShift+8, 360)
		return m, rainbowTick()
	case tea.KeyMsg:
		if m.screen == "input" {
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "esc", "b":
				m.screen = "menu"
				m.errorMsg = ""
				m.result = ""
			case "enter":
				path := strings.TrimSpace(m.input)
				if path == "" {
					m.errorMsg = "请输入文件路径，例如 test.txt"
					m.result = ""
					m.screen = "result"
					return m, nil
				}

				content, err := readFileContent(path)
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
			case "backspace":
				if len(m.input) > 0 {
					runes := []rune(m.input)
					m.input = string(runes[:len(runes)-1])
				}
			default:
				if len(msg.Runes) > 0 {
					m.input += string(msg.Runes)
				}
			}

			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "esc", "b":
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
			if m.screen == "menu" && m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			if m.screen != "menu" {
				return m, nil
			}

			choice := m.items[m.cursor]
			if choice == "退出" {
				m.quitting = true
				return m, tea.Quit
			}

			m.errorMsg = ""
			m.result = ""

			switch choice {
			case "读取文件（readFile）":
				m.screen = "input"
			}
		}
	}

	return m, nil
}

func renderColorLogo(hueShift float64) string {
	lines := strings.Split(strings.Trim(appLogo, "\n"), "\n")
	rendered := make([]string, 0, len(lines))

	for row, line := range lines {
		var builder strings.Builder
		for col, r := range []rune(line) {
			if r == ' ' {
				builder.WriteRune(r)
				continue
			}

			hue := math.Mod(hueShift+float64(row*20)+float64(col*6), 360)
			red, green, blue := hsvToRGB(hue, 0.9, 1)
			color := lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", red, green, blue))
			builder.WriteString(lipgloss.NewStyle().Foreground(color).Render(string(r)))
		}
		rendered = append(rendered, builder.String())
	}

	return strings.Join(rendered, "\n")
}

func hsvToRGB(h, s, v float64) (int, int, int) {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r1, g1, b1 float64
	switch {
	case h < 60:
		r1, g1, b1 = c, x, 0
	case h < 120:
		r1, g1, b1 = x, c, 0
	case h < 180:
		r1, g1, b1 = 0, c, x
	case h < 240:
		r1, g1, b1 = 0, x, c
	case h < 300:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}

	red := int(math.Round((r1 + m) * 255))
	green := int(math.Round((g1 + m) * 255))
	blue := int(math.Round((b1 + m) * 255))
	return red, green, blue
}

func (m lessonModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)

	if m.quitting {
		return "\n下次见！\n"
	}

	view := "\n" + renderColorLogo(m.hueShift) + "\n"
	view += titleStyle.Render("MOCLI--Made By MochiZen") + "\n\n"

	if m.screen == "menu" {
		for i, item := range m.items {
			cursor := " "
			if i == m.cursor {
				cursor = cursorStyle.Render(">")
			}
			view += fmt.Sprintf("%s %s\n", cursor, item)
		}
		view += "\n" + hintStyle.Render("j/k 或 方向键移动 · enter 执行 · q 退出") + "\n"
		return view
	}

	if m.screen == "input" {
		view += successStyle.Render("请输入要读取的文件路径") + "\n\n"
		view += fmt.Sprintf("> %s\n\n", m.input)
		view += hintStyle.Render("直接输入路径 · backspace 删除 · enter 读取 · esc 返回") + "\n"
		return view
	}

	view += successStyle.Render("执行结果") + "\n\n"
	if m.errorMsg != "" {
		view += errorStyle.Render(m.errorMsg) + "\n"
	} else {
		view += m.result + "\n"
	}
	view += "\n" + hintStyle.Render("按 b 或 esc 返回菜单 · q 退出") + "\n"
	return view
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "运行 Bubble Tea 入门交互界面",
	Run: func(cmd *cobra.Command, args []string) {
		program := tea.NewProgram(initialLessonModel())
		if _, err := program.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "运行 TUI 失败: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
