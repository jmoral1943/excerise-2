package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// JSON array of all the imports with the file name
var xi []map[string][]string

// Scans the file and searches for imports
func scanFile(f string) {
	// grabs the absolute path on any machine
	path, err := filepath.Abs(f)
	if err != nil {
		fmt.Println(err)
	}

	// Run the command to find the imports
	cmd, err := exec.Command("go", "list", "-f", "{{.Imports}}", path).Output()

	if err != nil {
		fmt.Println(err)
	}

	// Trim the string to then split into a splice
	s := strings.Trim(string(cmd), "[ ]\n")
	str := strings.Split(s, " ")

	// send the name of the file and the list of imports
	toJson(filepath.Base(f), str)
}
var wg sync.WaitGroup
var mutex sync.Mutex

func toJson(f string, imp []string) {
	// lock the slice so that no other go routine can have access to it
	mutex.Lock()
	m := map[string][]string{
		f: imp,
	}
	xi = append(xi, m)
	// unlocks the slice
	mutex.Unlock()
	// Finishes a waitgroup
	wg.Done()
}

// Goes through all the files in a dir
func listFiles(f string) {
	// Walks through the whole dir and searches for the deepest folder
	err := filepath.Walk(f,
		func(path string, info os.FileInfo, err error) error {
			matched, err := regexp.MatchString(`.go$`, filepath.Base(path))
			// if the file is a go file then I can check for imports
			if matched {
				// Adds a delta to the waitgroup for the go routine
				wg.Add(1)
				// go routine for scanning the file and getting the imports
				go scanFile(path)
			}
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	// Blocks and waits for the rest of the go routines to finish before continuing
	wg.Wait()
}

// Writes the data to a JSON file
func writeToFile() error {
	// Marshalling indenting the slice of imports and file names
	js, err := json.MarshalIndent(xi, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	// Creating a json file
	file, err := os.Create("imports.json")
	if err != nil {
		return err
	}
	// Preventing a memory leak
	defer file.Close()

	// Writing to the json file with the json
	_, err = file.Write(js)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func main() {
	listFiles(os.Args[1])
	err := writeToFile()
	if err != nil {
		log.Fatal(err)
	}
}
