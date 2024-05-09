package counter

import (
	"github.com/mistweaverco/timetrack.sh/internal/database"
	"github.com/mistweaverco/timetrack.sh/internal/stopwatch"
	"github.com/mistweaverco/timetrack.sh/internal/utils"
)

type Counter struct {
	ProjectName string
	TaskName    string
	Date        string
	Stopwatch   stopwatch.Model
}

type GetCounterOptions struct {
	ProjectName string
	TaskName    string
	Date        string
	Counters    []Counter
}

func convertStopwatchToSeconds(sw stopwatch.Model) int {
	time := sw.Elapsed()
	sec := time.Seconds()
	return int(sec)
}

func SaveAllCounters(counters []Counter) {
	for _, c := range counters {
		database.SaveCounter(c.ProjectName, c.TaskName, c.Date, convertStopwatchToSeconds(c.Stopwatch))
	}
}

func GetCounter(options GetCounterOptions) *Counter {
	for _, c := range options.Counters {
		if c.ProjectName == options.ProjectName && c.TaskName == options.TaskName && c.Date == options.Date {
			return &c
		}
	}
	return nil
}

func GetAllCounters() []Counter {
	var result []Counter
	counters := database.GetAllCounters()
	for _, c := range counters {
		result = append(result, Counter{
			ProjectName: c.ProjectName,
			TaskName:    c.TaskName,
			Date:        utils.ConvertDatetimeToDate(c.Date),
			Stopwatch:   stopwatch.New(c.Seconds),
		})
	}
	return result
}

func NewCounter(projectName, taskName, date string, elapsed int) Counter {
	return Counter{
		ProjectName: projectName,
		TaskName:    taskName,
		Date:        date,
		Stopwatch:   stopwatch.New(elapsed),
	}
}