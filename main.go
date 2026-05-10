package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	openrouter "github.com/revrost/go-openrouter"
)

const (
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
	colorReset  = "\033[0m"
)

func termWidth() int {
	w := os.Getenv("COLUMNS")
	if w != "" {
		if n, err := strconv.Atoi(w); err == nil && n > 20 {
			return n
		}
	}
	return 60
}

func boxTop(title string, color string) string {
	w := termWidth()
	title = " " + title + " "
	fill := w - len(title) - 2
	if fill < 1 {
		fill = 1
	}
	return color + "┌─" + title + strings.Repeat("─", fill) + colorReset
}

func boxBottom() string {
	w := termWidth()
	return colorGray + "└" + strings.Repeat("─", w-1) + colorReset
}

func boxContent(line string) string {
	return colorGray + "│ " + colorReset + line
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Print("\033[2J\033[H")

	scanner := bufio.NewScanner(os.Stdin)

	client := openrouter.NewClient(os.Getenv("OPENROUTER_KEY"))
	model := os.Getenv("MODEL")

	var messages []openrouter.ChatCompletionMessage
	messages = append(messages, openrouter.UserMessage(buildSystemPrompt()))

	inputTTY := func() bool {
		fi, _ := os.Stdin.Stat()
		return fi.Mode()&os.ModeCharDevice != 0
	}()

	for {
		var message string

		fmt.Print(colorCyan + boxTop("You", colorCyan) + "\n")
		fmt.Print(colorCyan + "│ " + colorReset)

		if !scanner.Scan() {
			break
		}
		message = scanner.Text()

		if inputTTY {
			fmt.Print("\n")
		} else {
			fmt.Print(message + "\n")
		}

		fmt.Print(colorCyan + boxBottom() + colorReset + "\n")

		if !handleCommand(client, message) {
			break
		}

		if message == "list" || message == "list free" {
			continue
		}

		messages = append(messages, openrouter.UserMessage(message))

		depth := 0
		const maxDepth = 5

		for {
			stream, err := callLLM(client, messages, model)
			if err != nil {
				fmt.Printf("\nLLM error: %v\n", err)
				break
			}

			var fullResponse strings.Builder

			modelName := ""
			first, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {
				fmt.Printf("\nStream error: %v\n", err)
				break
			}
			if err == nil {
				modelName = first.Model
				fmt.Print(boxTop("Assistant ("+modelName+")", colorGreen) + "\n")
				if len(first.Choices) > 0 {
					token := first.Choices[0].Delta.Content
					fmt.Print(token)
					fullResponse.WriteString(token)
				}
			}

			for {
				response, err := stream.Recv()

				if errors.Is(err, io.EOF) {
					break
				}

				if err != nil {
					fmt.Printf("\nStream error: %v\n", err)
					break
				}

				if len(response.Choices) > 0 {
					token := response.Choices[0].Delta.Content

					fmt.Print(token)

					fullResponse.WriteString(token)
				}
			}

			stream.Close()
			fmt.Print("\n" + boxBottom() + "\n\n")

			content := strings.TrimSpace(fullResponse.String())

			messages = append(messages, openrouter.AssistantMessage(content))

			toolCall := parseToolCall(content)

			if toolCall != nil && toolCall.Tool != "" {
				depth++

				if depth >= maxDepth {
					fmt.Print(colorYellow + "(max tool depth reached)" + colorReset + "\n")
					break
				}

				toolCall.CallID = newToolCallID()

				result := callTool(*toolCall)

				messages = append(messages, openrouter.ToolMessage(toolCall.CallID, result))

				label := toolCall.Tool
				if arg := toolCall.Arguments.Filename; arg != "" {
					label += " (" + arg + ")"
				} else if arg := toolCall.Arguments.Command; arg != "" {
					label += " (" + arg + ")"
				} else if arg := toolCall.Arguments.Pattern; arg != "" {
					label += " (" + arg + ")"
				}
				fmt.Print(boxTop("Tool: "+label, colorYellow) + "\n")
				for _, line := range strings.Split(strings.TrimRight(result, "\n"), "\n") {
					fmt.Print(boxContent(line) + "\n")
				}
				fmt.Print(boxBottom() + "\n\n")
			} else {
				break
			}
		}
	}

	fmt.Println(colorGray + "Thanks for using codegen" + colorReset)
}
