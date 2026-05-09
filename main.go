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
			fmt.Println("\n-------------------------------------------------\n")

			stream.Close()

			content := strings.TrimSpace(fullResponse.String())

			messages = append(messages, openrouter.AssistantMessage(content))

			var toolCall ToolCall

			err = json.Unmarshal([]byte(content), &toolCall)

			if err == nil && toolCall.Tool != "" {
				depth++

				if depth >= maxDepth {
					fmt.Println("\n\n(max tool depth reached)")
					break
				}

				fmt.Println()
				result := callTool(toolCall)

				messages = append(messages, openrouter.UserMessage("Tool result: "+result))

				fmt.Print("User: (tool result)\n")
			} else {
				break
			}
		}
	}

	fmt.Println("Thanks for using codegen")
}
