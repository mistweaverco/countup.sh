package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/mistweaverco/timetrack.sh/internal/counter"
	"github.com/mistweaverco/timetrack.sh/internal/database"
	"github.com/mistweaverco/timetrack.sh/internal/utils"
)

type ProjectTaskList struct {
	ProjectName string
	Task        []database.Project
	Table       table.Model
	TUIModel    *TUIModel
	KeyMap      KeyMap
}

func (m *ProjectTaskList) onWindowSizeMessage(msg tea.WindowSizeMsg) tea.Msg {
	// h, v := m.Table.GetFrameSize()
	// m.Table.Width(msg.Width - h)
	// m.Table.Height(msg.Height - v)
	return nil
}

func (m ProjectTaskList) Init() tea.Cmd {
	return nil
}

func (m ProjectTaskList) updateTaskCounters(msg tea.Msg) {
	for _, row := range m.Table.Rows() {
		taskName := row[0]
		taskDate := row[1]
		taskCounter := counter.GetCounter(counter.GetCounterOptions{
			Counters:    m.TUIModel.Counters,
			ProjectName: m.ProjectName,
			TaskName:    taskName,
			Date:        taskDate,
		})
		if taskCounter != nil {
			taskCounter.Stopwatch, _ = taskCounter.Stopwatch.Update(msg)
		}
	}
}

func (m ProjectTaskList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	selectedRow := m.Table.SelectedRow()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.onWindowSizeMessage(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Back):
			m.TUIModel.ActiveView = "projects"
			m.TUIModel.Views.ProjectsList = NewProjectsList(m.TUIModel)
			return m.TUIModel.Views.ProjectsList, nil
		case key.Matches(msg, m.KeyMap.CreateTask):
			m.TUIModel.Views.NewTaskForm = NewNewTaskForm(m.TUIModel, m.ProjectName)
			m.TUIModel.ActiveView = "newTask"
			m.TUIModel.Views.NewTaskForm.Init()
			return m.TUIModel.Views.NewTaskForm, nil
		case key.Matches(msg, m.KeyMap.EditTask):
			return m.TUIModel.Views.EditTaskForm, nil
		case key.Matches(msg, m.KeyMap.StartTask, m.KeyMap.StopTask):
			if selectedRow != nil {
				taskName := selectedRow[0]
				taskDate := selectedRow[1]
				taskCounter := counter.GetCounter(counter.GetCounterOptions{
					Counters:    m.TUIModel.Counters,
					ProjectName: m.ProjectName,
					TaskName:    taskName,
					Date:        taskDate,
				})
				if taskCounter != nil {
					sw := taskCounter.Stopwatch
					cmd := sw.Toggle()
					log.Info("TaskCounter", "Status", sw.Running())
					return m, cmd
				}
			}
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		}
	}
	m.updateTaskCounters(msg)
	m.Table, _ = m.Table.Update(msg)
	return m, tea.Batch(cmds...)
}

func (m ProjectTaskList) View() string {
	for idx, row := range m.Table.Rows() {
		seconds := row[2]
		counter := m.TUIModel.GetCounterForTask(m.ProjectName, row[0], row[1])
		if counter != nil {
			seconds = convertSecondsIntoTime(utils.FloatToInt(counter.Stopwatch.Elapsed().Seconds()))
		}
		row[2] = seconds
		m.Table.SetRows(append(m.Table.Rows()[:idx], append([]table.Row{row}, m.Table.Rows()[idx+1:]...)...))
	}
	return m.Table.View() + "\n"
}

func convertSecondsIntoTime(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds - (hours * 3600)) / 60
	seconds = seconds - (hours * 3600) - (minutes * 60)
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func convertDatetimeToDate(datetime string) string {
	return datetime[:10]
}

func NewProjectTaskList(tuiModel *TUIModel, projectName string) tea.Model {
	database.New()
	keyMap := DefaultKeyMap()
	m := ProjectTaskList{
		ProjectName: projectName,
		TUIModel:    tuiModel,
		KeyMap:      keyMap,
	}
	tasks := database.GetAllTasks(projectName)

	columns := []table.Column{
		{Title: projectName + " Tasks", Width: 30},
		{Title: "Date", Width: 30},
		{Title: "Time", Width: 20},
	}

	rows := []table.Row{}

	for _, task := range tasks {
		seconds := task.Seconds
		counter := tuiModel.GetCounterForTask(projectName, task.Name, task.Date)
		if counter != nil {
			seconds = utils.FloatToInt(counter.Stopwatch.Elapsed().Seconds())
		}
		rows = append(rows, []string{
			task.Name,
			convertDatetimeToDate(task.Date),
			convertSecondsIntoTime(seconds),
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