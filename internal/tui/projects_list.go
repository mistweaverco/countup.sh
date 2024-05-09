package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistweaverco/timetrack.sh/internal/database"
)

type ProjectsListModel struct {
	Projects []database.Project
	Table    table.Model
	TUIModel *TUIModel
	KeyMap   KeyMap
}

func (m *ProjectsListModel) onWindowSizeMessage(msg tea.WindowSizeMsg) tea.Msg {
	// h, v := m.Table.GetFrameSize()
	// m.Table.Width(msg.Width - h)
	// m.Table.Height(msg.Height - v)
	return nil
}

func (m *ProjectsListModel) deleteProject() {
	projectName := m.Table.SelectedRow()[0]
	project := database.Project{Name: projectName}
	database.DeleteProject(project)
	m.TUIModel.Views.ProjectsList = NewProjectsList(m.TUIModel)
}

func (m ProjectsListModel) Init() tea.Cmd {
	return nil
}

func (m ProjectsListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var projectName string
	selectedRow := m.Table.SelectedRow()
	if len(selectedRow) > 0 {
		projectName = selectedRow[0]
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.onWindowSizeMessage(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.ViewProject):
			m.TUIModel.ActiveView = "projectTaskList"
			m.TUIModel.Views.ProjectTaskList = NewProjectTaskList(m.TUIModel, projectName)
			return m.TUIModel.Views.ProjectTaskList, nil
		case key.Matches(msg, m.KeyMap.CreateProject):
			return m.TUIModel.Views.NewProjectForm, nil
		case key.Matches(msg, m.KeyMap.DeleteProject):
			m.deleteProject()
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		}
	}
	m.Table, _ = m.Table.Update(msg)
	return m, nil
}

func (m ProjectsListModel) View() string {
	return m.Table.View() + "\n"
}

func NewProjectsList(tuiModel *TUIModel) tea.Model {
	database.New()
	projects := database.GetAllProjects()
	keyMap := DefaultKeyMap()
	m := ProjectsListModel{
		Projects: projects,
		TUIModel: tuiModel,
		KeyMap:   keyMap,
	}
	columns := []table.Column{
		{Title: "Project", Width: 60},
		{Title: "Today's time", Width: 30},
	}

	rows := []table.Row{}

	for _, project := range m.Projects {
		rows = append(rows, []string{
			project.Name,
			"",
		})
	}

	m.Table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	return m
}