package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

const (
	AskQuestion = iota + 1
	AddFeature
	FixBug
	CodeCleanup
	WriteTests
	FindBugs
	CodeReview
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: ./gptswe [file1] [file2] ...")
		os.Exit(1)
	}

	// Get command line arguments
	fileList, commandFlag, detailsFlag := getCommandLineArguments()

	// Read files
	fileContents := readFiles(fileList)

	// Select action
	choice := commandFlag
	if choice == 0 {
		choice = selectAction()
	}

	// Build prompt
	prompt, err := buildPrompt(fileContents, choice, detailsFlag)
	handleError(err)

	fmt.Println(prompt)

	// Create job
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	// Create LLM Client
	llmClient, err := NewLLMClient()
	handleError(err)

	// Submit job
	_, err = llmClient.submitJob(choice, messages)
	handleError(err)
}

func commandDescriptions() map[int]string {
	return map[int]string{
		AskQuestion: "Ask a question",
		AddFeature:  "Add a feature",
		FixBug:      "Fix a bug",
		CodeCleanup: "Code cleanup",
		WriteTests:  "Write unit tests",
		FindBugs:    "Find bugs",
		CodeReview:  "Code review",
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
	}
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

func generateCommandDescriptionsText() string {
	// Dynamically generate command descriptions
	descText := "Skip the command choice prompt - "
	for i, cmd := range orderedCommands() {
		descText += fmt.Sprintf("%d: %s", cmd, commandDescriptions()[cmd])
		if i < len(orderedCommands())-1 {
			descText += ", "
		}
	}
	return descText
}

func commandCausesFileChanges(command int) bool {
	return command == AddFeature ||
		command == FixBug ||
		command == CodeCleanup ||
		command == WriteTests
}

func getCommandLineArguments() ([]string, int, string) {
	// Ensure arguments were provided
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		os.Exit(1)
	}

	// Get command line arguments
	var commandFlag int
	var detailsFlag string
	flag.IntVar(&commandFlag, "command", 0, generateCommandDescriptionsText())
	flag.StringVar(&detailsFlag, "details", "", "Provide additional details to the LLM")
	flag.Parse()

	// Use the remaining arguments as the list of files
	fileList := flag.Args()

	// Ensure files were provided
	if len(fileList) == 0 {
		printHelp()
		os.Exit(1)
	}

	return fileList, commandFlag, detailsFlag
}

func printHelp() {
	fmt.Println("Usage: ./gptswe [file1] [file2] ...")
	os.Exit(1)
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
