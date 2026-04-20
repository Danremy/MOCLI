package components

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type FileExplorer struct {
	table   table.Model
	cwd     string
	entries []fileEntry
}

type fileEntry struct {
	name  string
	path  string
	isDir bool
}

func NewFileExplorer(startPath string) FileExplorer {
	absPath, _ := filepath.Abs(startPath)
	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Name", Width: 48},
			{Title: "Type", Width: 8},
		}),
		table.WithFocused(true),
		table.WithHeight(16),
	)
	fe := FileExplorer{table: t, cwd: absPath}
	fe.reload()
	return fe
}

func (f *FileExplorer) reload() {
	files, err := os.ReadDir(f.cwd)
	if err != nil {
		return
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() != files[j].IsDir() {
			return files[i].IsDir()
		}
		return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
	})

	rows := make([]table.Row, len(files))
	f.entries = make([]fileEntry, len(files))
	for i, item := range files {
		typeText := "FILE"
		if item.IsDir() {
			typeText = "DIR"
		}
		f.entries[i] = fileEntry{
			name:  item.Name(),
			path:  filepath.Join(f.cwd, item.Name()),
			isDir: item.IsDir(),
		}
		rows[i] = table.Row{item.Name(), typeText}
	}
	f.table.SetRows(rows)
}

func (f *FileExplorer) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	f.table, cmd = f.table.Update(msg)
	return cmd
}

func (f *FileExplorer) View() string {
	return fmt.Sprintf("Current: %s\n\n%s", f.cwd, f.table.View())
}

func (f *FileExplorer) Selected() (string, bool) {
	cursor := f.table.Cursor()
	if cursor < 0 || cursor >= len(f.entries) {
		return "", false
	}
	entry := f.entries[cursor]
	if entry.isDir {
		f.cwd = entry.path
		f.reload()
		return "", false
	}
	return entry.path, true
}

func (f *FileExplorer) GoUp() {
	parent := filepath.Dir(f.cwd)
	if parent != f.cwd {
		f.cwd = parent
		f.reload()
	}
}
