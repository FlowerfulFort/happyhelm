package tui

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	pickvalues "github.com/flowerfulfort/happyhelm/internal/values"
)

var ErrCanceled = errors.New("selection canceled")

type model struct {
	entries  []pickvalues.ValueEntry
	query    string
	cursor   int
	offset   int
	width    int
	height   int
	selected map[int]bool
	canceled bool
	done     bool
}

func Pick(entries []pickvalues.ValueEntry, keywords []string) ([]pickvalues.ValueEntry, error) {
	m := model{
		entries:  entries,
		query:    strings.Join(keywords, " "),
		selected: make(map[int]bool),
	}

	tty, err := openTTY()
	if err != nil {
		return nil, err
	}
	defer tty.Close()

	program := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))
	result, err := program.Run()
	if err != nil {
		return nil, fmt.Errorf("run picker: %w", err)
	}

	final := result.(model)
	if final.canceled {
		return nil, ErrCanceled
	}

	var selected []pickvalues.ValueEntry
	for i, entry := range final.entries {
		if final.selected[i] {
			selected = append(selected, entry)
		}
	}
	return selected, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.clampOffset()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.canceled = true
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case "pgup", "ctrl+u":
			m.movePage(-1)
		case "pgdown", "ctrl+d":
			m.movePage(1)
		case " ", "space":
			if len(m.entries) > 0 {
				m.selected[m.cursor] = !m.selected[m.cursor]
			}
		case "a":
			m.toggleAll()
		}
	}
	m.clampOffset()
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Search: %s\n\n", m.query)

	start, end := m.visibleRange()
	for i := start; i < end; i++ {
		entry := m.entries[i]
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		check := " "
		if m.selected[i] {
			check = "x"
		}
		fmt.Fprintf(&b, "%s [%s] %-40s %v\n", cursor, check, entry.Path, entry.Value)
	}

	if start > 0 || end < len(m.entries) {
		fmt.Fprintf(&b, "\nshowing %d-%d of %d\n", start+1, end, len(m.entries))
	}
	b.WriteString("\nup/down or j/k: move    pgup/pgdown or ctrl+u/d: page    space: toggle    a: all    enter: confirm    q/esc: quit\n")
	return b.String()
}

func (m *model) toggleAll() {
	allSelected := len(m.entries) > 0
	for i := range m.entries {
		if !m.selected[i] {
			allSelected = false
			break
		}
	}

	if allSelected {
		m.selected = make(map[int]bool)
		return
	}
	for i := range m.entries {
		m.selected[i] = true
	}
}

func (m *model) movePage(direction int) {
	if len(m.entries) == 0 {
		return
	}

	delta := m.pageSize()
	if delta < 1 {
		delta = 1
	}

	m.cursor += direction * delta
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.entries) {
		m.cursor = len(m.entries) - 1
	}
}

func (m model) visibleRange() (int, int) {
	pageSize := m.pageSize()
	start := m.offset
	end := start + pageSize
	if end > len(m.entries) {
		end = len(m.entries)
	}
	return start, end
}

func (m model) pageSize() int {
	if m.height <= 0 {
		return len(m.entries)
	}
	size := m.height - 5
	if size < 1 {
		return 1
	}
	return size
}

func (m *model) clampOffset() {
	pageSize := m.pageSize()
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+pageSize {
		m.offset = m.cursor - pageSize + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}
