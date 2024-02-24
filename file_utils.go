package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// readFile reads the contents of the specified file and returns them as a slice of strings.
func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	var contents []string
	for {
		var line string
		_, err := fmt.Fscanln(file, &line)
		if err != nil {
			break
		}
		contents = append(contents, line)
	}

	return contents, nil
}

// createFile creates a new file with the specified contents.
func createFile(filePath string, contents string) error {
	return updateFile(filePath, contents)
}

// updateFile overwrites the given file with the provided contents.
func updateFile(filePath string, contents string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	for i, line := range strings.Split(contents, "\n") {
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

// runFunction parses the changes from the LLM response and applies them to the filesystem
func runFunction(function *openai.FunctionCall) (string, error) {
	// Args are a JSON string. Parse them into a map
	args := make(map[string]interface{})
	err := json.Unmarshal([]byte(function.Arguments), &args)
	if err != nil {
		return "", fmt.Errorf("failed to parse function args: %v", err)
	}

	switch function.Name {
	case "create_file":
		fileName, ok := args["file_name"]
		if !ok {
			return "", fmt.Errorf("missing file_name argument")
		}

		contents, ok := args["contents"]
		if !ok {
			return "", fmt.Errorf("missing contents argument")
		}

		// Ensure string types
		fileNameStr, ok := fileName.(string)
		if !ok {
			return "", fmt.Errorf("file_name argument must be a string")
		}

		contentsStr, ok := contents.(string)
		if !ok {
			return "", fmt.Errorf("contents argument must be a string")
		}

		err := createFile(fileNameStr, contentsStr)
		if err != nil {
			return "", err
		}

		return "", nil
	case "read_file":
		fileName, ok := args["file_name"]
		if !ok {
			return "", fmt.Errorf("missing file_name argument")
		}

		// Ensure string type
		fileNameStr, ok := fileName.(string)
		if !ok {
			return "", fmt.Errorf("file_name argument must be a string")
		}

		contents, err := readFile(fileNameStr)
		if err != nil {
			return "", err
		}

		return contents[0], nil
	case "update_file":
		fileName, ok := args["file_name"]
		if !ok {
			return "", fmt.Errorf("missing file_name argument")
		}

		contents, ok := args["contents"]
		if !ok {
			return "", fmt.Errorf("missing contents argument")
		}

		// Ensure string types
		fileNameStr, ok := fileName.(string)
		if !ok {
			return "", fmt.Errorf("file_name argument must be a string")
		}

		contentsStr, ok := contents.(string)
		if !ok {
			return "", fmt.Errorf("contents argument must be a string")
		}

		err := updateFile(fileNameStr, contentsStr)
		if err != nil {
			return "", err
		}

		return "", nil
	case "delete_file":
		fileName, ok := args["file_name"]
		if !ok {
			return "", fmt.Errorf("missing file_name argument")
		}

		// Ensure string type
		fileNameStr, ok := fileName.(string)
		if !ok {
			return "", fmt.Errorf("file_name argument must be a string")
		}

		err := deleteFile(fileNameStr)
		if err != nil {
			return "", err
		}

	case "finish":
		os.Exit(0)

	default:
		return "", fmt.Errorf("unknown command: %s", function.Name)
	}

	return "", nil
}
