package main

import (
	"fmt"
	"os"

	"github.com/rdbell/gptswe/logger"
	"github.com/sashabaranov/go-openai"
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
