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

	for {
		fmt.Print("User: ")

		scanner.Scan()

		message := scanner.Text()

		if message == "exit" {
			break
		}

		stream, err := callLLM(client, message)

		if err != nil {
			log.Fatal(err)
		}

		var fullResponse strings.Builder

		for {
			response, err := stream.Recv()

			if errors.Is(err, io.EOF) {
				fmt.Println("\nstream finished")
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

		stream.Close()

		content := strings.TrimSpace(fullResponse.String())

		var toolCall ToolCall

		err = json.Unmarshal([]byte(content), &toolCall)

		if err == nil && toolCall.Tool != "" {
			fmt.Println("\n\nTOOL DETECTED")
			toolRes := callTool(toolCall)
			fmt.Println(toolRes[0])

		}
	}

	fmt.Println("Thanks for using codegen")
}
