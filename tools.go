package main

import (
	"encoding/json"
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

	case "lsDir":
		return lsDir(toolCall.Arguments.Path)

	case "glob":
		return glob(toolCall.Arguments.Pattern)

	default:
		fmt.Println("unknown tool")
	}

	return ""
}

func readFile(filePath string) string {
	path := filepath.Clean(filePath)
	if !filepath.IsLocal(path) {
		return "readFile error: path must be local"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("readFile error: %v", err)
	}
	return string(data)
}

func writeFile(filePath string, content string) string {
	path := filepath.Clean(filePath)
	if !filepath.IsLocal(path) {
		return "writeFile error: path must be local"
	}
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Sprintf("writeFile error: %v", err)
	}
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

func parseToolCall(content string) *ToolCall {
	content = strings.TrimSpace(content)

	var toolCall ToolCall
	if err := json.Unmarshal([]byte(content), &toolCall); err == nil && toolCall.Tool != "" {
		return &toolCall
	}

	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && end > start {
		jsonStr := content[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &toolCall); err == nil && toolCall.Tool != "" {
			return &toolCall
		}
	}

	return nil
}

func lsDir(dirPath string) string {
	if dirPath == "" {
		dirPath = "."
	}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Sprintf("lsDir error: %v", err)
	}
	var b strings.Builder
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		b.WriteString(name)
		b.WriteString("\n")
	}
	return b.String()
}

func glob(pattern string) string {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Sprintf("glob error: %v", err)
	}
	if len(matches) == 0 {
		return "(no matches)"
	}
	return strings.Join(matches, "\n")
}
