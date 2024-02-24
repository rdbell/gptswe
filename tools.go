package main

import (
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func allTools() []openai.Tool {
	// List all files
	list := openai.FunctionDefinition{
		Name:        "list_files",
		Description: "List all files in the current directory",
	}

	// Create a new file
	create := openai.FunctionDefinition{
		Name:        "create_file",
		Description: "Create a new file",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"name": {
					Type:        jsonschema.String,
					Description: `The name of the file to create."`,
				},
				"contents": {
					Type:        jsonschema.String,
					Description: `The contents of the file."`,
				},
			},
		},
	}

	// Read a file
	read := openai.FunctionDefinition{
		Name:        "read_file",
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
		Name:        "update_file",
		Description: "Update a file",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"name": {
					Type:        jsonschema.String,
					Description: `The name of the file to update."`,
				},
				"contents": {
					Type:        jsonschema.String,
					Description: `The new contents of the file."`,
				},
			},
		},
	}

	// Delete a file
	del := openai.FunctionDefinition{
		Name:        "delete_file",
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

	l := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &list,
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

	return []openai.Tool{l, c, r, u, d, f}
}
