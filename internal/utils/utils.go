package utils

import (
	"database/sql"
	"embed"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DBFile string = "countup.sh.db"

var DBInstance *sql.DB

var ps = string(os.PathSeparator)

func GetDataDirectory() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	fullpath := dir + ps + "countup.sh" + ps
	mkdirerr := os.MkdirAll(fullpath, os.ModePerm)
	if mkdirerr != nil {
		log.Fatal(mkdirerr)
	}
	return fullpath
}

func DBNew() {
	userConfigDir := GetDataDirectory()
	fullPath := userConfigDir + DBFile
	_, staterr := os.Stat(fullPath)

	instance, err := sql.Open("sqlite3", fullPath)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	if DBInstance == nil {
		DBInstance = instance
	}
	if os.IsNotExist(staterr) {
		createAndPrefillDatabase()
	}
}

func DBClose() {
	DBInstance.Close()
}

//go:embed db.sql
var sqlfile embed.FS

func createAndPrefillDatabase() {
	readfile, err := sqlfile.ReadFile("db.sql")
	if err != nil {
		log.Fatal("Error reading db.sql: ", err)
	}
	str := string(readfile)
	_, err = DBInstance.Exec(str)
	if err != nil {
		log.Fatal("Error creating database: ", err, DBInstance.Stats())
	}
}

type DBCounter struct {
	Name    string
	Elapsed int
}

func DBGetCounter(name string) DBCounter {
	var counter DBCounter
	err := DBInstance.QueryRow("SELECT name, elapsed FROM counters WHERE name = ?", name).Scan(&counter.Name, &counter.Elapsed)
	if err != nil {
		counter = DBCounter{
			Name:    name,
			Elapsed: 0,
		}
	}
	return counter
}

func DBSetCounter(name string, elapsed int) {
	_, err := DBInstance.Exec("INSERT OR REPLACE INTO counters (name, elapsed) VALUES (?, ?)", name, elapsed)
	if err != nil {
		log.Fatal("Error setting counter: ", err)
	}
}