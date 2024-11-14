// main.go
package main

import (
	"fmt"
	"log"
	"perplexity-mini/services"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Prompt the user for a search query
	var query string
	fmt.Print("Enter search query: ")
	fmt.Scanln(&query)

	// Call the search service
	results, err := services.QueryGoogleSearch(query)
	if err != nil {
		log.Fatalf("Error querying Google Search: %v", err)
	}

	// Display the search results
	for _, item := range results {
		fmt.Printf("Title: %s\nLink: %s\nSnippet: %s\n\n", item.Title, item.Link, item.Snippet)
	}
}
