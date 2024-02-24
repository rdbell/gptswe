package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/rdbell/gptswe/logger"
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
			logger.Response(msg.Content)
		}

		if len(msg.ToolCalls) < 1 {
			logger.Error(errors.New("no tool selected"))
			dialogue = append(dialogue, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: "error: please select a tool to use, or use the 'finish' tool to complete the conversation",
			})

			continue
		}

		// Run the functions
		for _, toolCall := range msg.ToolCalls {
			logger.Tool(toolCall.Function.Name, toolCall.Function.Arguments)
			response := toolCall.Function
			out, err := runFunction(&response)
			if err != nil {
				logger.Error(err)
				out = fmt.Sprintf("Error: %v", err)
			}

			// Add the function result to the dialogue
			dialogue = append(dialogue, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    out,
				Name:       toolCall.Function.Name,
				ToolCallID: toolCall.ID,
			})
		}
	}
}
