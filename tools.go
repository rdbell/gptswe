package main

import (
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func allTools() []openai.Tool {
	// Create a new file
	create := openai.FunctionDefinition{
		Name:        "file_create",
		Description: "Create a new file",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"name": {
					Type:        jsonschema.String,
					Description: `The name of the file to create."`,
				},
				"content": {
					Type:        jsonschema.String,
					Description: `The contents of the file."`,
				},
			},
		},
	}

	// Read a file
	read := openai.FunctionDefinition{
		Name:        "file_read",
		Description: "Read a file",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"name": {
					Type:        jsonschema.String,
					Description: `The name of the file to read."`,
				},
			},
		},
	}

	// Update a file
	update := openai.FunctionDefinition{
		Name:        "file_update",
		Description: "Update a file",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"name": {
					Type:        jsonschema.String,
					Description: `The name of the file to update."`,
				},
				"content": {
					Type:        jsonschema.String,
					Description: `The new contents of the file."`,
				},
			},
		},
	}

	// Delete a file
	del := openai.FunctionDefinition{
		Name:        "file_delete",
		Description: "Delete a file",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"name": {
					Type:        jsonschema.String,
					Description: `The name of the file to delete."`,
				},
			},
		},
	}

	// Finish
	finish := openai.FunctionDefinition{
		Name:        "finish",
		Description: "Finish the conversation",
	}

	c := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &create,
	}

	r := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &read,
	}

	u := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &update,
	}

	d := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &del,
	}

	f := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &finish,
	}

	return []openai.Tool{c, r, u, d, f}
}
