package main

import (
	"bufio"
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

		if !handleCommand(client, message) {
			break
		}

		if message == "list" || message == "list free" {
			continue
		}

		messages = append(messages, openrouter.UserMessage(message))

		fmt.Println("\n-------------------------------------------------")

		depth := 0
		const maxDepth = 5

		for {
			stream, err := callLLM(client, messages)
			if err != nil {
				log.Fatal(err)
			}

			var fullResponse strings.Builder

			first, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {
				log.Fatal(err)
			}
			if err == nil {
				fmt.Printf("LLM name: %s\n", first.Model)
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
