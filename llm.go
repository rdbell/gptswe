package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type LLMClient struct {
	apiKey string
}

func NewLLMClient() (*LLMClient, error) {
	// Ensure that the API key is set.
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is not set")
	}

	return &LLMClient{apiKey: apiKey}, nil
}

// submitJob submits the messages to the LLM and receives the response.
func (client *LLMClient) submitJob(command int, messages []openai.ChatCompletionMessage) (string, error) {
	openAIClient := openai.NewClient(client.apiKey)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := openAIClient.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT4,
			Messages: messages,
			Stream:   true,
		},
	)
	if err != nil {
		return "", fmt.Errorf("ChatCompletionStream error: %v\n", err)
	}
	defer stream.Close()

	var out string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			break
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break
		}

		out += response.Choices[0].Delta.Content
		fmt.Print(response.Choices[0].Delta.Content)
		if len(out) >= 7 &&
			commandCausesFileChanges(command) &&
			!responseLooksValid(command, out) {
			return "", errors.New("response looks invalid")
		}
	}

	return out, nil
}

func responseLooksValid(command int, response string) bool {
	return strings.Index(response, "CREATE ") == 0 ||
		strings.Index(response, "UPDATE ") == 0 ||
		strings.Index(response, "DELETE ") == 0
}
