package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistweaverco/timetrack.sh/internal/database"
)

type NewProjectForm struct {
	Name     textinput.Model
	TUIModel *TUIModel
}

func NewNewProjectForm(tuiModel *TUIModel) tea.Model {
	form := NewProjectForm{}
	form.TUIModel = tuiModel
	form.Name = textinput.New()
	form.Name.Focus()
	return form
}

func (m *NewProjectForm) insertProject() {
	database.InsertProject(database.Project{Name: m.Name.Value()})
	m.TUIModel.ActiveView = "projects"
	m.TUIModel.Views.ProjectsList = NewProjectsList(m.TUIModel)
}

func (m NewProjectForm) Init() tea.Cmd {
	return nil
}

func (m NewProjectForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.Name.Value() != "" {
				m.insertProject()
				return m.TUIModel.Views.ProjectsList, nil
			} else {
				return m, nil
			}
		}
	}
	m.Name, cmd = m.Name.Update(msg)
	return m, cmd
}

func (m NewProjectForm) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.Name.View())
}