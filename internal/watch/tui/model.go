package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/bagaking/bilink/internal/output"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Status        string
	Added         int
	Removed       int
	Asking        bool
	ConfigSummary string

	frame      int
	showConfig bool
}

func NewModel(payload output.WatchPayload, configSummary string) Model {
	return Model{
		Status:        "watching",
		Added:         len(payload.Added),
		Removed:       len(payload.Removed),
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

const tickInterval = 250 * time.Millisecond

var askFrames = []string{"", ".", "..", "..."}

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
	var b strings.Builder
	b.WriteString("Bilink Watch\n")
	b.WriteString(fmt.Sprintf("status: %s\n", m.Status))
	b.WriteString(fmt.Sprintf("added: %d\nremoved: %d\n", m.Added, m.Removed))
	if m.Asking {
		b.WriteString(fmt.Sprintf("ASK%s  (y=accept, c=config, q=quit)\n", askFrames[m.frame]))
	} else {
		b.WriteString("press q to quit\n")
	}
	if m.showConfig && m.ConfigSummary != "" {
		b.WriteString("\nconfig:\n")
		b.WriteString(m.ConfigSummary)
		if !strings.HasSuffix(m.ConfigSummary, "\n") {
			b.WriteString("\n")
		}
	}
	return b.String()
}
