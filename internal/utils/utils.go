package utils

import (
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var ps = string(os.PathSeparator)

func GetDataDirectory() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	fullpath := dir + ps + "timetrack.sh" + ps
	mkdirerr := os.MkdirAll(fullpath, os.ModePerm)
	if mkdirerr != nil {
		log.Fatal(mkdirerr)
	}
	return fullpath
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func FloatToInt(input_num float64) int {
	return int(input_num)
}

func IntToString(input_num int) string {
	return strconv.Itoa(input_num)
}

func StringToInt(input_num string) int {
	i, err := strconv.Atoi(input_num)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func ConvertDatetimeToDate(datetime string) string {
	return datetime[:10]
}