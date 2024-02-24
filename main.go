package main

import (
	"fmt"
	"os"

	"github.com/rdbell/gptswe/logger"
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
	choice := selectAction()

	// Build prompt
	prompt, err := buildPrompt(choice)
	handleError(err)

	logger.Request(prompt)

	// Add to the dialoge
	dialogue = append(dialogue, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	// Create LLM Client
	llmClient, err := NewLLMClient()
	handleError(err)

	// Submit job
	err = llmClient.submitJob(dialogue)
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

func commandCausesFileChanges(command int) bool {
	return command == AddFeature ||
		command == FixBug ||
		command == CodeCleanup ||
		command == WriteTests
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
