package main

import (
	"log"
	"os"
	"strings"
)

func gotoStartDirectory() {
	evilproxy := "evilproxy.go"
	file, err := os.Open(evilproxy)
	if os.IsNotExist(err) {
		dir, err := os.Open("../")
		if err != nil {
			log.Fatalf("Unable to find %v directory. %v\n", evilproxy, err)
		}
		defer dir.Close()
		err = dir.Chdir()
		if err != nil {
			log.Fatalf("Unable to change directory %v\n", err)
		}
	}
	file.Close()
}

func getCurrentSubDirs() []string {
	dir, err := os.Open("./")
	if err != nil {
		log.Fatalf("Unable to open directory for reading. %v\n", err)
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		log.Fatalf("Unable to read directory. %v\n", err)
	}
	dirNames := make([]string, len(files))

	j := 0
	for i := 0; i != len(files); i++ {
		if files[i].IsDir() {
			if !strings.HasPrefix(files[i].Name(), ".") {
				dirNames[j] = "./" + files[i].Name()
				j++
			}
		}
	}

	return dirNames[0:j]
}
