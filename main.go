package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your search query: ")
	scanner.Scan()
	query = scanner.Text()

	// Get search results from Google API
	results, err := services.QueryGoogleSearch(query)
	if err != nil {
		log.Fatalf("Error querying Google Search: %v", err)
	}

	// // Collect URLs from search results
	// var urls []string
	// for _, item := range results {
	// 	urls = append(urls, item.Link)
	// }

	// Fetch content from URLs concurrently
	fmt.Println("Fetching content from URLs...")
	// fetchResults := services.ConcurrentFetchURLs(urls)

	// Print fetched content
	// for _, result := range fetchResults {
	// 	if result.Error != nil {
	// 		fmt.Printf("Error fetching %s: %v\n", result.URL, result.Error)
	// 	} else {
	// 		fmt.Printf("Fetched URL: %s\nContent: %s\n\n", result.URL, result.Body[:100]) // Print first 100 characters
	// 	}
	// }

	for _, items := range results {
		fmt.Println("Title:" + items.Title)
		fmt.Println("Link:" + items.Link)
		fmt.Println("Snippet:" + items.Snippet)
	}

	fmt.Println("collating all the responses")
	collatedResults := services.CollateSnippets(results)

	fmt.Println("Sending a request to OpenAI api")
	result, err := services.CallOpenAI(collatedResults)
	if err != nil {
		log.Fatalf("Error querying OPENAI API: %v", err)
	}
	fmt.Println("SUMMARY:" + result)

}
