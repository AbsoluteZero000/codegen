package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	openrouter "github.com/revrost/go-openrouter"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	scanner := bufio.NewScanner(os.Stdin)

	client := openrouter.NewClient(
		os.Getenv("OPENROUTER_KEY"),
		openrouter.WithXTitle("My App"),
		openrouter.WithHTTPReferer("https://myapp.com"),
	)

	for true {
		fmt.Printf("User: ")
		scanner.Scan()
		message := scanner.Text()

		if message == "exit" {
			break
		}

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openrouter.ChatCompletionRequest{
				Messages: []openrouter.ChatCompletionMessage{
					openrouter.UserMessage(buildPrompt(message)),
				},
			},
		)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return
		}
		fmt.Println(resp.Model+": ", resp.Choices[0].Message.Content)
	}

	fmt.Printf("Thanks for using codegen")

}
