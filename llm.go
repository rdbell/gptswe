package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var dialogue []openai.ChatCompletionMessage

func init() {
	dialogue = []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: "You are a software developer AI assistant. You are helping a user with their software project. " +
				"The user will ask you to perform various tasks related to their project. " +
				"You will need to understand the user's requests and use the provided tools to complete their reqests.",
		},
	}
}

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
func (client *LLMClient) submitJob(command int) (string, error) {
	openAIClient := openai.NewClient(client.apiKey)
	/*
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		stream, err := openAIClient.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model:    openai.GPT4,
				Messages: dialog,
				Tools:    allTools(),
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

			if len(response.Choices) == 0 {
				fmt.Println("\nNo choices in response")
				break
			}

			if len(response.Choices[

			out += response.Choices[0].Delta.Content
			fmt.Print(response.Choices[0].Delta.Content)
			if len(out) >= 7 &&
				commandCausesFileChanges(command) &&
				!responseLooksValid(command, out) {
				return "", errors.New("response looks invalid")
			}
		}
	*/

	// Call OpenAI
	resp, err := openAIClient.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    dialogue,
			Tools:       allTools(),
			Temperature: 1.2,
		},
	)

	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion error: %v\n", err)
	}

	// Get the response
	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) != 1 {
		fmt.Printf("Completion error: len(toolcalls): %v\n", len(msg.ToolCalls))
		return "", errors.New("no tool calls in response")
	}

	// Append the response to the dialogue
	dialogue = append(dialogue, msg)

	// Get the function call
	response := msg.ToolCalls[0].Function.Arguments

	// TODO: Validate the response

	// TODO: Call the function

	// TODO: Return the result

	return out, nil
}

func responseLooksValid(command int, response string) bool {
	return strings.Index(response, "CREATE ") == 0 ||
		strings.Index(response, "UPDATE ") == 0 ||
		strings.Index(response, "DELETE ") == 0
}
