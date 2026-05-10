package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	openrouter "github.com/revrost/go-openrouter"
)

func handleCommand(client *openrouter.Client, message string) bool {
	switch message {
	case "exit":
		return false
	case "list":
		models, err := client.ListModels(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nAvailable models:")
		for _, m := range models {
			fmt.Printf("  %s — %s\n", m.ID, m.Name)
		}
		fmt.Println()
		return true
	case "list free":
		models, err := client.ListModels(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nFree models:")
		for _, m := range models {
			if isFree(m) {
				fmt.Printf("  %s — %s\n", m.ID, m.Name)
			}
		}
		fmt.Println()
		return true
	default:
		return true
	}
}

func isFree(m openrouter.Model) bool {
	if strings.HasSuffix(m.ID, ":free") {
		return true
	}
	return m.Pricing.Prompt == "0" && m.Pricing.Completion == "0"
}
