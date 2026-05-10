package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	openrouter "github.com/revrost/go-openrouter"
)

func handleCommand(client *openrouter.Client, message string) bool {
	switch message {
	case "exit":
		return false
	case "list":
		printModels(client, "Available models:", false)
		return true
	case "list free":
		printModels(client, "Free models:", true)
		return true
	default:
		return true
	}
}

func printModels(client *openrouter.Client, header string, freeOnly bool) {
	models, err := client.ListModels(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing models: %v\n", err)
		return
	}
	fmt.Println("\n" + header)
	for _, m := range models {
		if freeOnly && !isFree(m) {
			continue
		}
		fmt.Printf("  %s — %s\n", m.ID, m.Name)
	}
	fmt.Println()
}

func isFree(m openrouter.Model) bool {
	if strings.HasSuffix(m.ID, ":free") {
		return true
	}
	return m.Pricing.Prompt == "0" && m.Pricing.Completion == "0"
}
