package main

import (
	"context"

	openrouter "github.com/revrost/go-openrouter"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func callLLM(client *openrouter.Client, messages []openrouter.ChatCompletionMessage, model string) (*openrouter.ChatCompletionStream, error) {
	return client.CreateChatCompletionStream(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
			Stream:   true,
		},
	)
}
