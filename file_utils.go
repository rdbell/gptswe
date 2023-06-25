package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// updateFile overwrites the given file with the provided contents.
func updateFile(filePath string, contents []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	for i, line := range contents {
		_, err := file.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write contents to file %s: %v", filePath, err)
		}

		// Add a new line for all lines except for the last one
		if i < len(contents)-1 {
			_, err := file.WriteString("\n")
			if err != nil {
				return fmt.Errorf("failed to write new line to file %s: %v", filePath, err)
			}
		}
	}

	return nil
}

// deleteFile deletes the specified file from the filesystem.
func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %v", filePath, err)
	}
	return nil
}

// applyChanges parses the changes from the LLM response and applies them to the filesystem
func applyChanges(changes string) error {
	// Select the first choice
	lines := strings.Split(changes, "\n")

	// Parse the response
	var fileStartIndex []int
	for i, line := range lines {
		if line == "FILE_START" {
			fileStartIndex = append(fileStartIndex, i+1)
		}
	}
	var fileEndIndex []int
	for i, line := range lines {
		if line == "FILE_END" {
			fileEndIndex = append(fileEndIndex, i)
		}
	}
	if len(fileStartIndex) != len(fileEndIndex) {
		return errors.New("fileStartIndex and fileEndIndex are unequal lengths")
	}
	var operations []string
	var fileNames []string
	var files [][]string
	for i := range fileStartIndex {
		ll := lines[fileStartIndex[i]:fileEndIndex[i]]
		files = append(files, ll)

		operationFileName := lines[fileStartIndex[i]-2]
		split := strings.Split(operationFileName, " ")
		if len(split) != 2 {
			return fmt.Errorf("encountered invalid operation/filename: %s", operationFileName)
		}

		operations = append(operations, split[0])
		fileNames = append(fileNames, split[1])

	}

	// Show the user a preview
	// Don't just print the raw LLM output - this is a better representation of what will actually happen on the disk
	for i, file := range files {
		fmt.Println(operations[i] + " " + fileNames[i])
		fmt.Println("FILE_START")
		for _, line := range file {
			fmt.Println(line)
		}
		fmt.Println("FILE_END")
		fmt.Println()
	}

	if os.Getenv("APPLY_CHANGES_NO_CONFIRM") != "true" {
		fmt.Println("Do you wish to apply these changes? (y/N)")
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		cont := strings.TrimSpace(scanner.Text())
		if cont != "y" && cont != "Y" {
			fmt.Println("Aborting")
			return nil
		}
	}

	// TODO: If this is not in a git repo, or if the branch has uncommitted changes, warn the user and prompt again

	for i, file := range files {
		var err error
		switch operations[i] {
		case "CREATE":
			err = updateFile(fileNames[i], file)
		case "UPDATE":
			err = updateFile(fileNames[i], file)
		case "DELETE":
			err = deleteFile(fileNames[i])
		default:
			return fmt.Errorf("encountered invalid operation %s", operations[i])
		}
		if err != nil {
			return err
		}

	}
	return nil
}
