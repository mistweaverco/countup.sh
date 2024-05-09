package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_ "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistweaverco/timetrack.sh/internal/counter"
	"github.com/mistweaverco/timetrack.sh/internal/database"
	"github.com/mistweaverco/timetrack.sh/internal/styles"
)

type KeyMap struct {
	LineUp              key.Binding
	LineDown            key.Binding
	Back                key.Binding
	CreateTask          key.Binding
	DeleteTask          key.Binding
	StartTask           key.Binding
	StopTask            key.Binding
	ViewProject         key.Binding
	CreateProject       key.Binding
	SubmitCreateProject key.Binding
	DeleteProject       key.Binding
	EditTask            key.Binding
	ShowFullHelp        key.Binding
	CloseFullHelp       key.Binding
	Quit                key.Binding
}

type TUIViews struct {
	ProjectsList    tea.Model
	NewProjectForm  tea.Model
	ProjectTaskList tea.Model
	NewTaskForm     tea.Model
	EditTaskForm    tea.Model
}

type TUIModel struct {
	Counters   []counter.Counter
	Styles     styles.Styles
	KeyMap     KeyMap
	ActiveView string
	Views      *TUIViews
	Help       help.Model
}

func (m *TUIModel) onWindowSizeMessage(msg tea.WindowSizeMsg) tea.Msg {
	// h, v := columnStyle.GetFrameSize()
	// columnStyle.Width(msg.Width - h)
	// columnStyle.Height(msg.Height - v)
	return nil
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.LineUp, km.LineDown, km.CreateProject, km.ShowFullHelp, km.Quit}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.LineUp, km.LineDown, km.CreateProject, km.Quit},
		{km.DeleteProject, km.CreateTask, km.DeleteTask, km.ViewProject},
		{km.StartTask, km.ViewProject, km.CloseFullHelp},
	}
}

func DefaultKeyMap() KeyMap {
	const spacebar = " "
	return KeyMap{
		Back: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "go back"),
		),
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		ViewProject: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view project tasks"),
		),
		CreateProject: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "create new project"),
		),
		DeleteProject: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "delete project"),
		),
		EditTask: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "edit task"),
		),
		CreateTask: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "create new task"),
		),
		StartTask: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "start/continue a task"),
		),
		StopTask: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "stop/pause a task"),
		),
		DeleteTask: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "delete task"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),
	}
}

func (m TUIModel) ShortHelp() []key.Binding {
	return m.KeyMap.ShortHelp()
}

func (m TUIModel) FullHelp() [][]key.Binding {
	return m.KeyMap.FullHelp()
}

func (m *TUIModel) updateKeybindings() {
	if m.Help.ShowAll {
		m.KeyMap.ShowFullHelp.SetEnabled(true)
		m.KeyMap.CloseFullHelp.SetEnabled(true)
	} else {
		m.KeyMap.ShowFullHelp.SetEnabled(true)
		m.KeyMap.CloseFullHelp.SetEnabled(false)
	}
	switch m.ActiveView {
	case "projects":
		m.KeyMap.ViewProject.SetEnabled(true)
		m.KeyMap.CreateProject.SetEnabled(true)
		m.KeyMap.DeleteProject.SetEnabled(true)
		m.KeyMap.CreateTask.SetEnabled(false)
		m.KeyMap.DeleteTask.SetEnabled(false)
		m.KeyMap.StartTask.SetEnabled(false)
	case "newProject":
		m.KeyMap.ViewProject.SetEnabled(false)
		m.KeyMap.SubmitCreateProject.SetEnabled(true)
		m.KeyMap.LineUp.SetEnabled(false)
		m.KeyMap.LineDown.SetEnabled(false)
		m.KeyMap.DeleteProject.SetEnabled(false)
		m.KeyMap.CreateTask.SetEnabled(false)
		m.KeyMap.DeleteTask.SetEnabled(false)
		m.KeyMap.StartTask.SetEnabled(false)
	}
}

func (m TUIModel) Init() tea.Cmd {
	return nil
}

func (m TUIModel) helpView() string {
	return m.Styles.HelpStyle.Render(m.Help.View(m))
}

func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.updateKeybindings()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.onWindowSizeMessage(msg)
	}

	switch m.ActiveView {
	case "projects":
		m.Views.ProjectsList, cmd = m.Views.ProjectsList.Update(msg)
		return m, cmd
	case "newProject":
		m.Views.NewProjectForm, cmd = m.Views.NewProjectForm.Update(msg)
		return m, cmd
	case "projectTaskList":
		m.Views.ProjectTaskList, cmd = m.Views.ProjectTaskList.Update(msg)
		return m, cmd
	case "newTask":
		m.Views.NewTaskForm, cmd = m.Views.NewTaskForm.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m TUIModel) View() string {
	switch m.ActiveView {
	case "projects":
		return m.Views.ProjectsList.View() + "\n" + m.helpView()
	case "newProject":
		return m.Views.NewProjectForm.View() + "\n" + m.helpView()
	case "projectTaskList":
		return m.Views.ProjectTaskList.View() + "\n" + m.helpView()
	case "newTask":
		return m.Views.NewTaskForm.View() + "\n" + m.helpView()
	default:
		return ""
	}
}

func (m *TUIModel) GetCounterForTask(projectName, taskName, date string) *counter.Counter {
	for _, c := range m.Counters {
		if c.ProjectName == projectName && c.TaskName == taskName && c.Date == date {
			return &c
		}
	}
	return nil
}

func Start() {
	database.New()
	m := TUIModel{
		ActiveView: "projects",
		Views:      &TUIViews{},
		KeyMap:     DefaultKeyMap(),
		Styles:     styles.DefaultStyles(),
		Help:       help.New(),
		Counters:   counter.GetAllCounters(),
	}
	m.Views.ProjectsList = NewProjectsList(&m)
	m.Views.NewProjectForm = NewNewProjectForm(&m)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}