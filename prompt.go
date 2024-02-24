package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// buildPrompt prompts the user for more details based on their chosen command.
func buildPrompt(command int) (string, error) {
	// Build instructions
	instructions := "I have a software project in my current directory. I need help with the following: \n\n"

	switch command {
	case AskQuestion:
		fmt.Println("Ask a question about your code:")
	case AddFeature:
		instructions += "Implement the following feature(s): "
		fmt.Println("Describe the feature:")
	case FixBug:
		instructions += "Fix the following bug(s): "
		fmt.Println("Describe the bug:")
	case CodeCleanup:
		instructions += "Cleanup and refactor this codebase. "
		fmt.Println("Additional comments:")
	case WriteTests:
		instructions += "Write unit tests for this codebase. "
		fmt.Println("Additional comments:")
	case FindBugs:
		instructions += "Find bugs in this codebase. "
		fmt.Println("Additional comments:")
	case CodeReview:
		instructions += "Make suggestions for code improvements. "
		fmt.Println("Additional comments:")
	default:
		return "", errors.New("invalid selection")
	}

	instructions += "\n\n"

	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if strings.TrimSpace(scanner.Text()) != "" {
		instructions += strings.TrimSpace(scanner.Text())
	}

	instructions += "\n\n" + "Interact with my files using the provided functions. First fetch the file list, then create, read, update, or delete files as needed."

	return instructions, nil
}
