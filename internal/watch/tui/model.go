package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/bagaking/bilink/internal/output"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Status        string
	Added         []string
	Removed       []string
	Asking        bool
	ConfigSummary string

	frame      int
	showConfig bool
}

const (
	tickInterval = 250 * time.Millisecond
	maxChanges   = 8
)

var askFrames = []string{"", ".", "..", "..."}

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7DF9FF"))
	panelStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#7D5FFF")).Padding(0, 1)
	accentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6AD5"))
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7DF9FF"))
)

func NewModel(payload output.WatchPayload, configSummary string) Model {
	return Model{
		Status:        "watching",
		Added:         payload.Added,
		Removed:       payload.Removed,
		Asking:        true,
		ConfigSummary: configSummary,
	}
}

func Run(payload output.WatchPayload, configSummary string) error {
	program := tea.NewProgram(NewModel(payload, configSummary), tea.WithAltScreen())
	_, err := program.Run()
	return err
}

type tickMsg struct{}

func (m Model) Init() tea.Cmd {
	return tea.Tick(tickInterval, func(time.Time) tea.Msg { return tickMsg{} })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.Asking {
			m.frame = (m.frame + 1) % len(askFrames)
			return m, tea.Tick(tickInterval, func(time.Time) tea.Msg { return tickMsg{} })
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "c":
			m.showConfig = !m.showConfig
			return m, nil
		case "y", "enter":
			m.Asking = false
			m.Status = "confirmed"
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	title := titleStyle.Render("BILINK WATCH")
	stats := panelStyle.Render(strings.Join([]string{
		accentStyle.Render(fmt.Sprintf("STATUS: %s", strings.ToUpper(m.Status))),
		fmt.Sprintf("ADDED: %d", len(m.Added)),
		fmt.Sprintf("REMOVED: %d", len(m.Removed)),
	}, "\n"))

	changes := panelStyle.Render(renderChanges(m.changeLines()))
	body := lipgloss.JoinHorizontal(lipgloss.Top, stats, "  ", changes)

	footer := footerStyle.Render(fmt.Sprintf("ASK%s  :help :config  y=accept  q=quit", askFrames[m.frame]))
	if !m.Asking {
		footer = mutedStyle.Render("press q to quit")
	}

	sections := []string{title, body}
	if m.showConfig && m.ConfigSummary != "" {
		sections = append(sections, panelStyle.Render("CONFIG\n"+strings.TrimRight(m.ConfigSummary, "\n")))
	}
	sections = append(sections, footer)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) changeLines() []string {
	var lines []string
	for _, path := range m.Added {
		lines = append(lines, "+ "+path)
	}
	for _, path := range m.Removed {
		lines = append(lines, "- "+path)
	}
	return tail(lines, maxChanges)
}

func renderChanges(lines []string) string {
	if len(lines) == 0 {
		return mutedStyle.Render("no changes")
	}
	return strings.Join(lines, "\n")
}

func tail(items []string, max int) []string {
	if len(items) <= max {
		return items
	}
	return items[len(items)-max:]
}
