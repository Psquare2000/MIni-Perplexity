package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// CallOpenAI sends a request to OpenAI's API to summarize content into 100 words or less
func CallOpenAI(content string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not set in environment variables")
	}

	// Prepare the OpenAI API request payload
	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo", // You can switch to gpt-4 or other models if required
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant that summarizes content into 200 words or more"},
			{"role": "user", "content": fmt.Sprintf("Please summarize the following content, pretend you're perplexity ai:\n\n%s", content)},
		},
		"max_tokens":  150, // Allow enough room for 100 words
		"temperature": 0.7, // Adjust for creativity vs. determinism
	}

	// Serialize the payload to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to serialize request: %v", err)
	}

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the API call
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check for API errors
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		return "", fmt.Errorf("API error: %v", errorResponse)
	}

	// Parse the response
	var openAIResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract the summary
	if len(openAIResponse.Choices) > 0 {
		return openAIResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no content in response")
}
