// services/search_service.go
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const googleSearchAPI = "https://www.googleapis.com/customsearch/v1"

// SearchResult represents each individual search result
type SearchResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

// GoogleResponse represents the response structure from the Google Custom Search API
type GoogleResponse struct {
	Items []SearchResult `json:"items"`
}

// QueryGoogleSearch performs a search query using Google Custom Search API
func QueryGoogleSearch(query string) ([]SearchResult, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")           // Get API key from environment variables
	searchEngineID := os.Getenv("SEARCH_ENGINE_ID") // Get Search Engine ID from environment variables

	if apiKey == "" || searchEngineID == "" {
		return nil, fmt.Errorf("API key or Search Engine ID not set")
	}

	// Build the request URL with query, API key, and search engine ID
	requestURL := fmt.Sprintf("%s?key=%s&cx=%s&q=%s",
		googleSearchAPI, apiKey, searchEngineID, url.QueryEscape(query))

	// Make the HTTP GET request
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search results: %v", err)
	}
	defer resp.Body.Close()

	// Parse the JSON response into GoogleResponse struct
	var googleResponse GoogleResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Return the search results
	return googleResponse.Items, nil
}

func CollateSnippets(results []SearchResult) string {
	var snippets []string

	// Loop through each result and collect snippets
	for _, result := range results {
		if result.Snippet != "" {
			snippets = append(snippets, result.Snippet)
		}
	}

	// Join all snippets with a newline or space for better readability
	return strings.Join(snippets, "\n")
}
