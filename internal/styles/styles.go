package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mistweaverco/timetrack.sh/internal/colors"
)

type Styles struct {
	TableHeader lipgloss.Style
	TitleBar    lipgloss.Style
	Title       lipgloss.Style
	HelpStyle   lipgloss.Style
}

func DefaultStyles() (s Styles) {
	c := colors.DefaultColors()
	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)
	s.TableHeader = lipgloss.NewStyle().Foreground(lipgloss.Color(c.Red)).Background(lipgloss.Color(c.Red)).Padding(0, 1)
	return s
}