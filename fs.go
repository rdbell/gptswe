package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// readFiles reads the given files and returns their contents in a format compatible with the prompt.
func readFiles(args []string) string {
	fileContents := &strings.Builder{}
	for _, file := range args {
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Error reading file %s: %s", file, err)
		} else {
			fmt.Fprintf(fileContents, "PRINT %s\n", file)
			fileContents.WriteString("FILE_START\n")
			fileContents.WriteString(string(data))
			fileContents.WriteString("\nFILE_END\n\n")
		}
	}

	return fileContents.String()
}
