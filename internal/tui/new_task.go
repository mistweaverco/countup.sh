package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistweaverco/timetrack.sh/internal/database"
)

type NewTaskForm struct {
	ProjectName string
	Form        *huh.Form
	TUIModel    *TUIModel
	Confirm     bool
}

func NewNewTaskForm(tuiModel *TUIModel, projectName string) NewTaskForm {
	m := NewTaskForm{
		ProjectName: projectName,
		Confirm:     false,
	}
	m.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Key("name").Title("Name of task").Prompt("?"),
			huh.NewText().Title("The notes / description of the task").Key("description"),
			huh.NewConfirm().Affirmative("Create Task").Negative("Cancel").Key("confirm").Value(&m.Confirm),
		),
	)
	m.TUIModel = tuiModel
	return m
}

func (m NewTaskForm) createTask(name string, description string) {
	database.InsertTask(database.Task{Name: name, Description: description, ProjectName: m.ProjectName})
}

func (m NewTaskForm) Init() tea.Cmd {
	m.Form.Run()
	return nil
}

func (m NewTaskForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if m.Form != nil {
		f, cmd := m.Form.Update(msg)
		m.Form = f.(*huh.Form)
		cmds = append(cmds, cmd)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	if m.Form.State == huh.StateCompleted {
		name := m.Form.GetString("name")
		description := m.Form.GetString("description")
		m.createTask(name, description)
		m.TUIModel.ActiveView = "projectTaskList"
		m.TUIModel.Views.ProjectTaskList = NewProjectTaskList(m.TUIModel, m.ProjectName)
		return m.TUIModel.Views.ProjectTaskList, nil
	}
	return m, tea.Batch(cmds...)
}

func (m NewTaskForm) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.Form.View())
}