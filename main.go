package main

import (
	"io/fs"
	"log"
	"os"
	"strings"
)

func main() {
	workingDirectory := os.Args[1]

	entries, err := os.ReadDir(workingDirectory)
	if err != nil {
		log.Fatalln(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		extensionStart := strings.LastIndex(entry.Name(), ".")
		extension := strings.ToLower(entry.Name()[extensionStart+1:])

		checkIfDirExist(workingDirectory, extension)

		targetDirectory := workingDirectory + "/" + extension

		if !fileExistsInTargetDirectory(targetDirectory, entry.Name()) {
			err = os.Rename(workingDirectory+"/"+entry.Name(), targetDirectory+"/"+entry.Name())
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func checkIfDirExist(workingDirectory, dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.Mkdir(workingDirectory+"/"+dirName, fs.FileMode(0700))
	} else {
		return
	}
}

func fileExistsInTargetDirectory(targetDirectory, fileName string) bool {
	if _, err := os.Stat(targetDirectory + "/" + fileName); os.IsExist(err) {
		return true
	} else {
		return false
	}
}
