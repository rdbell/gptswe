package logger

import (
	"github.com/fatih/color"
)

// Request logs a request to the LLM.
func Request(msg string) {
	// Display banner
	banner("request")

	// Display question
	message("Submitted query", msg, color.Cyan)
}

// Response logs a response from the LLM.
func Response(msg string) {
	// Display banner
	banner("response")

	// Display answer
	message("Received response", msg, color.Green)
}

// Tool logs a tool being used by the LLM.
func Tool(function string, args string, response string) {
	// Display banner
	banner("tool")

	// Display function and args
	message("Running function: ", function+" "+args, color.Yellow)
	message("Function response: ", response, color.Yellow)
}

// Error logs an error from the LLM.
func Error(err error) {
	// Display banner
	banner("error")

	// Display error
	message("Received error", err.Error(), color.Red)
}
