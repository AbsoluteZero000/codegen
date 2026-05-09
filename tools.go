package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func callTool(toolCall ToolCall) string {
	switch toolCall.Tool {

	case "readFile":
		result := readFile(toolCall.Arguments.Filename)
		return result

	default:
		fmt.Println("unknown tool")
	}

	return ""
}

func readFile(filePath string) string {
	wd, err := os.Getwd()
	check(err)

	path := filepath.Join(wd, filePath)
	data, err := os.ReadFile(path)
	check(err)

	return string(data)
}
