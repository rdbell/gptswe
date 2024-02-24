package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/pp"
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
func (client *LLMClient) submitJob(dialogue []openai.ChatCompletionMessage) error {
	openAIClient := openai.NewClient(client.apiKey)

	for {
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
			return fmt.Errorf("CreateChatCompletion error: %v\n", err)
		}

		// Get the response
		msg := resp.Choices[0].Message

		// Append the response to the dialogue
		dialogue = append(dialogue, msg)

		if msg.Content != "" {
			fmt.Println(msg.Content)
		}

		if len(msg.ToolCalls) != 1 {
			fmt.Println("No tool calls found")
			continue
		}

		pp.Println(msg.ToolCalls[0].Function)

		// Run the function
		response := msg.ToolCalls[0].Function
		out, err := runFunction(&response)
		if err != nil {
			fmt.Println("Error: ", err)
			out = fmt.Sprintf("Error: %v", err)
		}

		// Add the function result to the dialogue
		dialogue = append(dialogue, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    out,
			Name:       msg.ToolCalls[0].Function.Name,
			ToolCallID: msg.ToolCalls[0].ID,
		})
	}
}

func responseLooksValid(command int, response string) bool {
	return strings.Index(response, "CREATE ") == 0 ||
		strings.Index(response, "UPDATE ") == 0 ||
		strings.Index(response, "DELETE ") == 0
}
