package main

import (
	"encoding/json"
	"log"
	"os"
)

var ps = string(os.PathSeparator)

func getCurrentWorkingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getJsonRoot() JsonRoot {
	currentDir := getCurrentWorkingDirectory()
	file, err := os.Open(currentDir + ps + "packages.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var jsonRoot JsonRoot
	err = decoder.Decode(&jsonRoot)
	if err != nil {
		log.Fatal(err)
	}

	return jsonRoot
}