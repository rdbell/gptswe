package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// listFiles returns a directory tree of the current working directory.
func listFiles() (string, error) {
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// Ignore hidden files and directories by checking the entire path.
		// Split the path and check each part.
		for _, part := range strings.Split(path, string(os.PathSeparator)) {
			if strings.HasPrefix(part, ".") && part != "." && part != ".." {
				return filepath.SkipDir
			}
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to list files: %v", err)
	}

	return strings.Join(files, "\n"), nil
}

// readFile reads the contents of the specified file and returns them as a slice of strings.
func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	var contents []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contents = append(contents, scanner.Text())
	}

	return strings.Join(contents, "\n"), nil
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
	case "list_files":
		files, err := listFiles()
		if err != nil {
			return "", err
		}

		return files, nil
	case "create_file":
		fileName, ok := args["name"]
		if !ok {
			return "", fmt.Errorf("missing name argument")
		}

		contents, ok := args["contents"]
		if !ok {
			return "", fmt.Errorf("missing contents argument")
		}

		// Ensure string types
		fileNameStr, ok := fileName.(string)
		if !ok {
			return "", fmt.Errorf("name argument must be a string")
		}

		contentsStr, ok := contents.(string)
		if !ok {
			return "", fmt.Errorf("contents argument must be a string")
		}

		err := createFile(fileNameStr, contentsStr)
		if err != nil {
			return "", err
		}
	case "read_file":
		fileName, ok := args["name"]
		if !ok {
			return "", fmt.Errorf("missing name argument")
		}

		// Ensure string type
		fileNameStr, ok := fileName.(string)
		if !ok {
			return "", fmt.Errorf("name argument must be a string")
		}

		contents, err := readFile(fileNameStr)
		if err != nil {
			return "", err
		}

		return contents, nil
	case "update_file":
		fileName, ok := args["name"]
		if !ok {
			return "", fmt.Errorf("missing name argument")
		}

		contents, ok := args["contents"]
		if !ok {
			return "", fmt.Errorf("missing contents argument")
		}

		// Ensure string types
		fileNameStr, ok := fileName.(string)
		if !ok {
			return "", fmt.Errorf("name argument must be a string")
		}

		contentsStr, ok := contents.(string)
		if !ok {
			return "", fmt.Errorf("contents argument must be a string")
		}

		err := updateFile(fileNameStr, contentsStr)
		if err != nil {
			return "", err
		}
	case "delete_file":
		fileName, ok := args["name"]
		if !ok {
			return "", fmt.Errorf("missing name argument")
		}

		// Ensure string type
		fileNameStr, ok := fileName.(string)
		if !ok {
			return "", fmt.Errorf("name argument must be a string")
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

	return "Success", nil
}
