package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type fileEntry struct {
	name  string
	path  string
	isDir bool
	size  int64
}

type explorerModel struct {
	cwd          string
	table        table.Model
	entries      []fileEntry
	errMsg       string
	inPreview    bool
	previewPath  string
	previewText  string
	selectedPath string
	quitting     bool
}

func newExplorerModel(startPath string) explorerModel {
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		absPath = "."
	}

	columns := []table.Column{
		{Title: "Name", Width: 48},
		{Title: "Type", Width: 8},
		{Title: "Size", Width: 12},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(16),
	)

	m := explorerModel{
		cwd:   absPath,
		table: t,
	}
	m.reloadEntries()

	return m
}

func (m explorerModel) Init() tea.Cmd {
	return nil
}

func (m *explorerModel) reloadEntries() {
	files, err := os.ReadDir(m.cwd)
	if err != nil {
		m.errMsg = fmt.Sprintf("读取目录失败: %v", err)
		m.entries = nil
		m.table.SetRows([]table.Row{})
		return
	}

	sort.Slice(files, func(i, j int) bool {
		left, right := files[i], files[j]
		if left.IsDir() != right.IsDir() {
			return left.IsDir()
		}
		return strings.ToLower(left.Name()) < strings.ToLower(right.Name())
	})

	rows := make([]table.Row, 0, len(files))
	entries := make([]fileEntry, 0, len(files))

	for _, item := range files {
		info, infoErr := item.Info()
		sizeText := "-"
		entrySize := int64(0)
		if infoErr == nil {
			entrySize = info.Size()
			if !item.IsDir() {
				sizeText = fmt.Sprintf("%d", entrySize)
			}
		}

		typeText := "FILE"
		if item.IsDir() {
			typeText = "DIR"
		}

		entries = append(entries, fileEntry{
			name:  item.Name(),
			path:  filepath.Join(m.cwd, item.Name()),
			isDir: item.IsDir(),
			size:  entrySize,
		})
		rows = append(rows, table.Row{item.Name(), typeText, sizeText})
	}

	m.entries = entries
	m.errMsg = ""
	m.table.SetRows(rows)
	if len(rows) == 0 {
		m.table.SetCursor(0)
		return
	}

	if m.table.Cursor() >= len(rows) {
		m.table.SetCursor(len(rows) - 1)
	}
}

func (m explorerModel) selectedEntry() (fileEntry, bool) {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.entries) {
		return fileEntry{}, false
	}
	return m.entries[cursor], true
}

func (m *explorerModel) enterSelected() {
	entry, ok := m.selectedEntry()
	if !ok {
		return
	}

	if entry.isDir {
		m.cwd = entry.path
		m.reloadEntries()
		return
	}

	content, err := os.ReadFile(entry.path)
	if err != nil {
		m.errMsg = fmt.Sprintf("读取文件失败: %v", err)
		return
	}

	preview := string(content)
	if len(preview) > 4096 {
		preview = preview[:4096] + "\n\n... (已截断)"
	}

	if strings.TrimSpace(preview) == "" {
		preview = "（文件为空）"
	}

	m.inPreview = true
	m.previewPath = entry.path
	m.previewText = preview
	m.errMsg = ""
}

func (m *explorerModel) goParent() {
	parent := filepath.Dir(m.cwd)
	if parent == m.cwd {
		return
	}
	m.cwd = parent
	m.reloadEntries()
}

func (m explorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch key := msg.(type) {
	case tea.KeyMsg:
		switch key.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.inPreview {
				m.inPreview = false
				m.previewPath = ""
				m.previewText = ""
			}
			return m, nil
		case "backspace", "h":
			if m.inPreview {
				m.inPreview = false
				m.previewPath = ""
				m.previewText = ""
				return m, nil
			}
			m.goParent()
			return m, nil
		case "enter":
			if m.inPreview {
				m.selectedPath = m.previewPath
				return m, tea.Quit
			}

			entry, ok := m.selectedEntry()
			if !ok {
				return m, nil
			}

			if entry.isDir {
				m.cwd = entry.path
				m.reloadEntries()
				return m, nil
			}

			m.selectedPath = entry.path
			return m, tea.Quit
		case "l", "p":
			if m.inPreview {
				return m, nil
			}
			m.enterSelected()
			return m, nil
		}
	}

	if !m.inPreview {
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m explorerModel) View() string {
	if m.quitting {
		return "\nBye!\n"
	}

	if m.inPreview {
		return fmt.Sprintf(
			"\nPreview: %s\n\n%s\n\nesc/backspace 返回列表 · q 退出\n",
			m.previewPath,
			m.previewText,
		)
	}

	view := fmt.Sprintf("\nRF Graphics Explorer\nCurrent: %s\n\n%s\n", m.cwd, m.table.View())
	if m.errMsg != "" {
		view += "\n" + m.errMsg + "\n"
	}
	view += "\nj/k 或方向键移动 · enter 选择文件/进入目录 · l 或 p 预览文件 · backspace/h 返回上级 · q 退出\n"
	return view
}

func RunFileExplorer(startPath string) (string, error) {
	// Run in alternate screen to avoid leaving TUI help lines in the main terminal.
	program := tea.NewProgram(newExplorerModel(startPath), tea.WithAltScreen())
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}

	explorer, ok := finalModel.(explorerModel)
	if !ok {
		return "", fmt.Errorf("unexpected explorer model type: %T", finalModel)
	}

	return explorer.selectedPath, nil
}
