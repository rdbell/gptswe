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
		instructions += "\n\n" +
			"Interact with files using the following commands: CREATE, UPDATE, DELETE\n\n" +
			"Reply in the following format:\n\n" +
			"```\n" +
			"UPDATE file1.ext\n" +
			"FILE_START\n" +
			"updated file contents...\n" +
			"FILE_END\n" +
			"UPDATE file2.ext\n" +
			"FILE_START\n" +
			"updated file contents...\n" +
			"FILE_END\n" +
			"DELETE filename3.ext\n" +
			"FILE_START\n" +
			"(deleted)\n" +
			"FILE_END\n" +
			"CREATE filename3.ext\n" +
			"FILE_START\n" +
			"new file contents...\n" +
			"FILE_END\n" +
			"(etc.)\n" +
			"```\n\n" +
			"YOUR RESPONSE SHOULD STRICTLY ADHERE TO THE ABOVE FORMAT.\n" +
			"- Only include files and commands relevant to the task at hand.\n" +
			"- Do not include explanations.\n" +
			"- Do not abbreviate or omit file contents - include the entire file contents in your response.\n" +
			"- Ensure to keep all existing comments within the code while making changes.\n" +
			"- DO NOT include anything in your response other than commands to interact with the files.\n" +
			"- DO NOT skip lines by writing comments like `// ommiting existing functions...` or `// code is unchanged` or ` // ... `."
	}

	return instructions, nil
}
