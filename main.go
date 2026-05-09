package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	openrouter "github.com/revrost/go-openrouter"
)

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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	scanner := bufio.NewScanner(os.Stdin)

	client := openrouter.NewClient(os.Getenv("OPENROUTER_KEY"))

	var messages []openrouter.ChatCompletionMessage
	messages = append(messages, openrouter.UserMessage(buildSystemPrompt()))

	for {
		fmt.Print("User: ")

		if !scanner.Scan() {
			break
		}

		message := scanner.Text()

		if message == "exit" {
			break
		}

		messages = append(messages, openrouter.UserMessage(message))

		depth := 0
		const maxDepth = 5

		for {
			stream, err := callLLM(client, messages)
			if err != nil {
				log.Fatal(err)
			}

			var fullResponse strings.Builder

			for {
				response, err := stream.Recv()

				if errors.Is(err, io.EOF) {
					break
				}

				if err != nil {
					log.Fatal(err)
				}

				if len(response.Choices) > 0 {
					token := response.Choices[0].Delta.Content

					fmt.Print(token)

					fullResponse.WriteString(token)
				}
			}
			fmt.Println("\n-------------------------------------------------")

			stream.Close()

			content := strings.TrimSpace(fullResponse.String())

			messages = append(messages, openrouter.AssistantMessage(content))

			toolCall := parseToolCall(content)

			if toolCall != nil && toolCall.Tool != "" {
				depth++

				if depth >= maxDepth {
					fmt.Println("\n(max tool depth reached)")
					break
				}

				fmt.Println()
				result := callTool(*toolCall)

				messages = append(messages, openrouter.UserMessage("Tool result: "+result))

				fmt.Print("User: (tool result)\n")
			} else {
				break
			}
		}
	}

	fmt.Println("Thanks for using codegen")
}
