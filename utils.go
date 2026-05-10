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

func callLLM(client *openrouter.Client, messages []openrouter.ChatCompletionMessage) (*openrouter.ChatCompletionStream, error) {
	return client.CreateChatCompletionStream(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Messages: messages,
			Stream:   true,
		},
	)
}
