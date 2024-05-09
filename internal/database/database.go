package database

import (
	"database/sql"
	"embed"
	"os"

	"github.com/charmbracelet/log"
	"github.com/mistweaverco/timetrack.sh/internal/utils"
)

const DBFile string = "timetrack.sh.db"

var Instance *sql.DB

func New() {
	userConfigDir := utils.GetDataDirectory()
	fullPath := userConfigDir + DBFile
	_, staterr := os.Stat(fullPath)

	instance, err := sql.Open("sqlite3", fullPath)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	if Instance == nil {
		Instance = instance
	}
	if os.IsNotExist(staterr) {
		createAndPrefillDatabase()
	}
}

func Close() {
	Instance.Close()
}

//go:embed db.sql
var sqlfile embed.FS

func createAndPrefillDatabase() {
	readfile, err := sqlfile.ReadFile("db.sql")
	if err != nil {
		log.Fatal("Error reading db.sql: ", err)
	}
	str := string(readfile)
	_, err = Instance.Exec(str)
	if err != nil {
		log.Fatal("Error creating database: ", err, Instance.Stats())
	}
}

type Task struct {
	Name        string
	Description string
	Date        string
	Seconds     int
	ProjectName string
}

func GetAllTasks(projectName string) []Task {
	rows, err := Instance.Query("SELECT name, description, date, seconds FROM tasks WHERE project_name = ?", projectName)
	if err != nil {
		log.Fatal("Error getting tasks: ", err)
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Name, &task.Description, &task.Date, &task.Seconds)
		if err != nil {
			log.Fatal("Error scanning tasks: ", err)
		}
		tasks = append(tasks, Task{
			Name:        task.Name,
			Description: task.Description,
			Date:        task.Date,
			Seconds:     task.Seconds,
		})
	}
	return tasks
}

type Project struct {
	Name  string
	Tasks []Task
}

func DeleteProject(project Project) {
	_, err := Instance.Exec("DELETE FROM projects WHERE name = ?", project.Name)
	if err != nil {
		log.Fatal("Error deleting project: ", err)
	}
}

func InsertProject(project Project) {
	_, err := Instance.Exec("INSERT INTO projects (name) VALUES (?)", project.Name)
	if err != nil {
		log.Fatal("Error saving project: ", err)
	}
}

func InsertTask(task Task) {
	_, err := Instance.Exec("INSERT INTO tasks (name, description, project_name) VALUES (?, ?, ?)", task.Name, task.Description, task.ProjectName)
	if err != nil {
		log.Fatal("Error saving project: ", err)
	}
}

func GetAllProjects() []Project {
	rows, err := Instance.Query("SELECT name FROM projects")
	if err != nil {
		log.Fatal("Error getting projects: ", err)
	}
	defer rows.Close()
	var projects []Project
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal("Error scanning projects: ", err)
		}
		projects = append(projects, Project{
			Name: name,
		})
		rows_tasks, err_tasks := Instance.Query("SELECT name, description, date, seconds FROM tasks WHERE project_name = ?", name)
		if err_tasks != nil {
			log.Fatal("Error getting tasks: ", err)
		}
		defer rows_tasks.Close()
		var tasks []Task
		for rows_tasks.Next() {
			var task Task
			err = rows_tasks.Scan(&task.Name, &task.Description, &task.Date, &task.Seconds)
			if err != nil {
				log.Fatal("Error scanning tasks: ", err)
			}
			tasks = append(tasks, Task{
				Name:        task.Name,
				Description: task.Description,
				Date:        task.Date,
				Seconds:     task.Seconds,
			})
		}
		projects[len(projects)-1].Tasks = tasks
	}
	return projects
}

type Counter struct {
	Elapsed int
}

func GetCounter(projectName string, taskName string, date string) Counter {
	var counter Counter
	err := Instance.QueryRow("SELECT elapsed FROM tasks WHERE project_name = ? AND name = ? AND date = ?", projectName, taskName, date).Scan(&counter.Elapsed)
	if err != nil {
		counter = Counter{
			Elapsed: 0,
		}
	}
	return counter
}

type GetAllCountersResultItem struct {
	TaskName    string
	ProjectName string
	Date        string
	Seconds     int
}

func GetAllCounters() []GetAllCountersResultItem {
	counters := []GetAllCountersResultItem{}
	rows, err := Instance.Query("SELECT name, project_name, date, seconds FROM tasks ORDER BY date DESC")
	if err != nil {
		log.Fatal("Error getting counters: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		c := GetAllCountersResultItem{}
		err = rows.Scan(&c.TaskName, &c.ProjectName, &c.Date, &c.Seconds)
		if err != nil {
			log.Fatal("Error scanning counters: ", err)
		}
		counters = append(counters, c)
	}
	return counters
}

type GetAllTodayCountersResultItem struct {
	TaskName    string
	ProjectName string
	Date        string
	Seconds     int
}

func GetAllTodayCounters() []GetAllTodayCountersResultItem {
	rows, err := Instance.Query("SELECT name, project_name, date, seconds FROM tasks WHERE date = date('now')")
	if err != nil {
		log.Fatal("Error getting counters: ", err)
	}
	defer rows.Close()
	var counters []GetAllTodayCountersResultItem
	for rows.Next() {
		var counter GetAllTodayCountersResultItem
		err = rows.Scan(&counter.TaskName, &counter.ProjectName, &counter.Date, &counter.Seconds)
		if err != nil {
			log.Fatal("Error scanning counters: ", err)
		}
		counters = append(counters, GetAllTodayCountersResultItem{
			TaskName:    counter.TaskName,
			ProjectName: counter.ProjectName,
			Date:        counter.Date,
			Seconds:     counter.Seconds,
		})
	}
	return counters
}

func SaveCounter(projectName string, taskName string, date string, elapsed int) {
	_, err := Instance.Exec("INSERT OR REPLACE INTO tasks (seconds) VALUES (?) WHERE project_name = ? AND name = ? AND date = ?", elapsed, projectName, taskName, elapsed)
	if err != nil {
		log.Fatal("Error setting counter: ", err)
	}
}