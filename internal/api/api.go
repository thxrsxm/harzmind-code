// Package api handles communication with external Large Language Model APIs,
// including sending messages, retrieving available models, and managing API requests.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message represents a message in the chat.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a chat request.
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse represents a chat response.
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// ModelsResponse represents a models response.
type ModelsResponse struct {
	Data []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

// SendMessage sends a message to the API and returns the response.
// It takes the API URL, model, token, and messages as input.
// It returns the response content and an error if any.
func SendMessage(url, model, token string, messages []Message) (string, error) {
	// Check if the URL is valid.
	if len(url) == 0 {
		return "", fmt.Errorf("url is not valid")
	}
	// Check if the model is valid.
	if len(model) == 0 {
		return "", fmt.Errorf("model is not valid")
	}
	// Check if the token is valid.
	if len(token) == 0 {
		return "", fmt.Errorf("API token is not valid")
	}
	// Create a new chat request.
	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
	}
	// Marshal the request body to JSON.
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	// Create a new HTTP request.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	// Set the content type to JSON.
	req.Header.Set("Content-Type", "application/json")
	// Set the authorization header.
	req.Header.Set("Authorization", "Bearer "+token)
	// Add timeout to prevent hanging requests
	client := &http.Client{Timeout: 120 * time.Second}
	// Send the request and get the response.
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}
	// Unmarshal the response to a ChatResponse.
	var chatResp ChatResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	// Check if the response has any choices.
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return chatResp.Choices[0].Message.Content, nil
}

// GetModels gets the list of models from the API.
// It takes the API URL and token as input.
// It returns the list of models and an error if any.
func GetModels(api, token string) ([]string, error) {
	// Create the models URL.
	modelsURL := strings.TrimSuffix(api, "/chat/completions") + "/models"
	// Create a new HTTP request.
	req, err := http.NewRequest("GET", modelsURL, nil)
	if err != nil {
		return nil, err
	}
	// Set the authorization header.
	req.Header.Set("Authorization", "Bearer "+token)
	// Create a new HTTP client.
	client := &http.Client{}
	// Send the request and get the response.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API Error: %s - %s", resp.Status, string(body))
	}
	// Unmarshal the response to a ModelsResponse.
	var modelsResp ModelsResponse
	err = json.Unmarshal(body, &modelsResp)
	if err != nil {
		return nil, err
	}
	// Create a new list of models.
	models := []string{}
	for _, model := range modelsResp.Data {
		models = append(models, model.ID)
	}
	return models, nil
}
