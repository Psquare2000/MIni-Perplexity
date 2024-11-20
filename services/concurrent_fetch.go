package services

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// FetchContent fetches and extracts meaningful content from a URL
func FetchContent(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a GET request with a custom User-Agent
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for URL %s: %v", url, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Check for HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL %s: status code %d", url, resp.StatusCode)
	}

	// Parse the HTML with goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML for URL %s: %v", url, err)
	}

	// Extract meaningful content (e.g., paragraphs)
	content := ""
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		content += s.Text() + "\n"
	})

	// Clean up content (optional)
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return "", fmt.Errorf("no meaningful content found at URL %s", url)
	}

	if len(content) < 100 { // Skip minimal or placeholder content
		return "", fmt.Errorf("content too short or invalid at URL %s", url)
	}

	return content, nil
}

func ConcurrentFetchContent(urls []string) string {
	var wg sync.WaitGroup
	results := make(chan string, len(urls)) // Buffered channel for results

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			content, err := FetchContent(u)
			if err != nil {
				fmt.Printf("Error fetching %s: %v\n", u, err)
				return
			}
			results <- content
		}(url)
	}

	// Wait for all Goroutines to complete
	wg.Wait()
	close(results)

	// Aggregate valid content
	var aggregatedContent string
	for result := range results {
		aggregatedContent += result + "\n"
	}

	return aggregatedContent
}
