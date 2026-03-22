package tui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const appLogo = `
███╗   ███╗ ██████╗  ██████╗██╗     ██╗
████╗ ████║██╔═══██╗██╔════╝██║     ██║
██╔████╔██║██║   ██║██║     ██║     ██║
██║╚██╔╝██║██║   ██║██║     ██║     ██║
██║ ╚═╝ ██║╚██████╔╝╚██████╗███████╗██║
╚═╝     ╚═╝ ╚═════╝  ╚═════╝╚══════╝╚═╝
`

var (
	logoGradientStart = [3]int{147, 112, 219}
	logoGradientMid   = [3]int{138, 43, 226}
	logoGradientEnd   = [3]int{106, 90, 205}
)

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)

	if m.quitting {
		return "\n下次见！\n"
	}

	if m.screen == "result" {
		view := "\n" + titleStyle.Render("Read File Result") + "\n\n"
		if m.errorMsg != "" {
			view += errorStyle.Render(m.resultView.View()) + "\n"
		} else {
			view += m.resultView.View() + "\n"
		}
		view += "\n" + hintStyle.Render("j/k 或方向键滚动 · pgup/pgdown 翻页 · g/G 到顶部/底部 · esc 返回菜单 · q 退出") + "\n"
		return view
	}

	view := "\n" + renderColorLogo(m.hueShift) + "\n"
	view += titleStyle.Render("MOCLI--Made By MochiZen") + "\n\n"

	if m.screen == "menu" {
		for i, item := range m.choices {
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
		view += m.textInput.View() + "\n\n"
		view += hintStyle.Render("直接输入路径 · backspace 删除 · enter 读取 · esc 返回") + "\n"
		return view
	}

	view += successStyle.Render("执行结果") + "\n\n"
	if m.errorMsg != "" {
		view += errorStyle.Render(m.errorMsg) + "\n"
	} else {
		view += m.result + "\n"
	}
	view += "\n" + hintStyle.Render("按 esc 返回菜单 · q 退出") + "\n"
	return view
}

func renderColorLogo(hueShift float64) string {
	lines := strings.Split(strings.Trim(appLogo, "\n"), "\n")
	rendered := make([]string, 0, len(lines))

	for _, line := range lines {
		rendered = append(rendered, renderGradientLine(line, hueShift))
	}

	return strings.Join(rendered, "\n")
}

func renderGradientLine(line string, shift float64) string {
	var builder strings.Builder
	runes := []rune(line)

	for col, r := range runes {
		if r == ' ' {
			builder.WriteRune(r)
			continue
		}

		factor := 0.5 + 0.5*math.Sin(float64(col)*0.35-shift)
		color := bluePurpleColor(factor)
		builder.WriteString(lipgloss.NewStyle().Foreground(color).Render(string(r)))
	}

	return builder.String()
}

func bluePurpleColor(t float64) lipgloss.Color {
	var from, to [3]int
	localT := t
	if t < 0.5 {
		from = logoGradientStart
		to = logoGradientMid
		localT = t * 2
	} else {
		from = logoGradientMid
		to = logoGradientEnd
		localT = (t - 0.5) * 2
	}

	red := mixChannel(from[0], to[0], localT)
	green := mixChannel(from[1], to[1], localT)
	blue := mixChannel(from[2], to[2], localT)
	return lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", red, green, blue))
}

func mixChannel(start, end int, t float64) int {
	value := float64(start) + (float64(end)-float64(start))*t
	return int(math.Round(value))
}
