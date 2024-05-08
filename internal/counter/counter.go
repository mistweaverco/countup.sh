package counter

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistweaverco/countup.sh/internal/stopwatch"
	"github.com/mistweaverco/countup.sh/internal/utils"
)

const (
	timerNameColor = "#2EF8BB"
	stopwatchColor = "#FF5F87"
)

var (
	timerNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(timerNameColor)).MarginRight(1)
	stopwatchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(stopwatchColor)).MarginRight(1)
	helpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(2)
)

var baseTimerStyle = lipgloss.NewStyle().Padding(1, 2)

var (
	lastID int
	idMtx  sync.Mutex
)

type model struct {
	timerName string
	stopwatch stopwatch.Model
	keymap    keymap
	help      help.Model
	quitting  bool
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

func (m model) Init() tea.Cmd {
	return nil
}

func getFormattedTimeString(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func (m model) BackupDBPeriodic() {
	elapsed := int(m.stopwatch.Elapsed().Seconds())
	if elapsed%30 == 0 {
		m.UpdateElapsedInDB()
	}
}

func (m model) View() string {
	s := ""
	sw := getFormattedTimeString(m.stopwatch.Elapsed()) + "\n"
	if !m.quitting {
		m.BackupDBPeriodic()
		s = timerNameStyle.Render(m.timerName) + " " + sw
		s += m.helpView()
	} else {
		m.UpdateElapsedInDB()
	}
	return s
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			return m, m.stopwatch.Reset()
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			m.keymap.stop.SetEnabled(!m.stopwatch.Running())
			m.keymap.start.SetEnabled(m.stopwatch.Running())
			return m, m.stopwatch.Toggle()
		}
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m model) UpdateElapsedInDB() {
	elapsed := int(m.stopwatch.Elapsed().Seconds())
	utils.DBSetCounter(m.timerName, elapsed)
}

func Start(timerName string) {
	utils.DBNew()
	t := utils.DBGetCounter(timerName)
	m := model{
		timerName: timerName,
		stopwatch: stopwatch.NewWithInterval(time.Second, t.Elapsed),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
	m.keymap.stop.SetEnabled(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}