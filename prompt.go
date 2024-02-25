package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	AskQuestion = iota + 1
	AddFeature
	FixBug
	CodeCleanup
	WriteTests
	FindBugs
	CodeReview
	Lint
)

func commandDescriptions() map[int]string {
	return map[int]string{
		AskQuestion: "Ask a question",
		AddFeature:  "Add a feature",
		FixBug:      "Fix a bug",
		CodeCleanup: "Code cleanup",
		WriteTests:  "Write unit tests",
		FindBugs:    "Find bugs",
		CodeReview:  "Code review",
		Lint:        "Fix linter errors",
	}
}

func orderedCommands() []int {
	return []int{
		AskQuestion,
		AddFeature,
		FixBug,
		CodeCleanup,
		WriteTests,
		FindBugs,
		CodeReview,
		Lint,
	}
}

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
	case Lint:
		instructions += "Run the linter and fix any errors/warnings. "
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

func selectAction() int {
	// Print list of commands
	fmt.Println("Choose from the following commands:")
	for _, cmd := range orderedCommands() {
		fmt.Printf("%d: %s\n", cmd, commandDescriptions()[cmd])
	}
	fmt.Print("> ")

	// Read user's choice
	var choice int
	for {
		_, err := fmt.Scanln(&choice)

		if err == nil && commandDescriptions()[choice] != "" {
			break
		}

		fmt.Println("Please choose a valid command")
	}

	return choice
}
