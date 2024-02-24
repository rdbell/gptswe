package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// buildPrompt prompts the user for more details based on their chosen command.
func buildPrompt(fileContents string, command int, detailsFlag string) (string, error) {
	// Build instructions
	instructions := fmt.Sprintf("I have the following files in my project.\n\n%s\n\n", fileContents)

	if detailsFlag == "" {
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
	} else {
		instructions += detailsFlag
	}

	if commandCausesFileChanges(command) {
		instructions += "\n\n" + "Interact with files using the provided functions."
	}

	return instructions, nil
}
