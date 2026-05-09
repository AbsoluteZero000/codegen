package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func callTool(toolCall ToolCall) string {
	switch toolCall.Tool {

	case "readFile":
		return readFile(toolCall.Arguments.Filename)

	case "writeFile":
		return writeFile(toolCall.Arguments.Filename, toolCall.Arguments.Content)

	case "bash":
		return bash(toolCall.Arguments.Command)

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

func bash(command string) string {
	command = sanitize(command)
	if command == "" {
		return "bash error: blocked suspicious command"
	}

	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("bash error: %v\n%s", err, output)
	}
	return string(output)
}

func sanitize(command string) string {
	command = strings.TrimSpace(command)

	dangerous := []string{"rm -rf /", "rm -rf /*", "dd if=", ":(){:|:&};:", "mkfs", "fdisk", "parted"}
	lower := strings.ToLower(command)
	for _, d := range dangerous {
		if strings.Contains(lower, d) {
			return ""
		}
	}

	if strings.Contains(command, "..") {
		wd, err := os.Getwd()
		if err != nil {
			return ""
		}
		clean := filepath.Clean(filepath.FromSlash(command))
		if !strings.HasPrefix(clean, wd) {
			return ""
		}
	}

	return command
}
