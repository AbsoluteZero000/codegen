package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func callTool(toolCall ToolCall) string {
	switch toolCall.Tool {

	case "readFile":
		return readFile(toolCall.Arguments.Filename)

	case "writeFile":
		return writeFile(toolCall.Arguments.Filename, toolCall.Arguments.Content)

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

func writeFile(filePath string, content string) string {
	wd, err := os.Getwd()
	check(err)

	path := filepath.Join(wd, filePath)
	err = os.WriteFile(path, []byte(content), 0644)
	check(err)

	return "File written: " + path
}
