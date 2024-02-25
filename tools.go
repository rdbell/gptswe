package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rdbell/gptswe/logger"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func allTools() []openai.Tool {
	// Tool definitions simplified using the utility function
	tooldDefinitions := []struct {
		name        string
		description string
		parameters  map[string]jsonschema.Definition
	}{
		{
			name:        "list_files",
			description: "List all files in the current directory",
			parameters:  nil,
		},
		{
			name:        "create_file",
			description: "Create a new file",
			parameters: map[string]jsonschema.Definition{
				"name":     {Type: jsonschema.String, Description: `The name of the file to create."`},
				"contents": {Type: jsonschema.String, Description: `The contents of the file."`},
			},
		},
		{
			name:        "read_file",
			description: "Read a file",
			parameters: map[string]jsonschema.Definition{
				"name": {Type: jsonschema.String, Description: `The name of the file to read."`},
			},
		},
		{
			name:        "update_file",
			description: "Update a file",
			parameters: map[string]jsonschema.Definition{
				"name":     {Type: jsonschema.String, Description: `The name of the file to update."`},
				"contents": {Type: jsonschema.String, Description: `The new contents of the file."`},
			},
		},
		{
			name:        "delete_file",
			description: "Delete a file",
			parameters: map[string]jsonschema.Definition{
				"name": {Type: jsonschema.String, Description: `The name of the file to delete."`},
			},
		},
		{
			name:        "run_linter",
			description: "Run the code linter on the project",
			parameters:  nil,
		},
		{
			name:        "finish",
			description: "Finish the conversation",
			parameters:  nil,
		},
	}

	var tools []openai.Tool
	for _, def := range tooldDefinitions {
		tool := openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{Name: def.name, Description: def.description, Parameters: jsonschema.Definition{Type: jsonschema.Object, Properties: def.parameters}},
		}
		tools = append(tools, tool)
	}

	return tools
}

// listFiles returns a directory tree of the current working directory.
func listFiles() (string, error) {
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the path or any of its parent directories start with a dot.
		isHidden := false
		for _, part := range strings.Split(path, string(os.PathSeparator)) {
			if strings.HasPrefix(part, ".") && part != "." && part != ".." {
				isHidden = true
				break
			}
		}

		// Only add non-hidden files/directories to the list.
		if !isHidden {
			files = append(files, path)
		}
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

// runLinter runs the golangci-lint linter on the project.
func runLinter() (out string, err error) {
	// Define the command and arguments
	cmd := exec.Command("golangci-lint", "run", "./...")

	// Create buffers to capture STDOUT and STDERR
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run the command
	err = cmd.Run()

	// Combine STDOUT and STDERR
	out = stdoutBuf.String() + stderrBuf.String()

	// Return combined output and an error indicating the command failed
	if err != nil {
		if out == "" {
			out = err.Error()
		}

		return out, fmt.Errorf("linter returned the following: \n%s", out)
	}

	// If no error, just return the output
	return out, nil
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
	case "run_linter":
		out, err := runLinter()
		if err != nil {
			return "", err
		}

		return out, nil
	case "finish":
		logger.Tool(function.Name, function.Arguments, "Finished conversation")
		os.Exit(0)

	default:
		return "", fmt.Errorf("unknown command: %s", function.Name)
	}

	return "Success", nil
}
