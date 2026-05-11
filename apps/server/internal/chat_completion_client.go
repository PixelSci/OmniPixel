package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"omni-pixel/domain"
)

type ChatCompletionClient struct {
	httpClient *http.Client
}

func NewChatCompletionClient(timeout time.Duration) *ChatCompletionClient {
	return &ChatCompletionClient{
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (client *ChatCompletionClient) Complete(request domain.ChatCompletionRequest) (string, error) {
	provider := strings.TrimSpace(strings.ToLower(request.Provider))
	model := strings.TrimSpace(request.Model)
	apiKey := strings.TrimSpace(request.APIKey)

	if provider == "" || model == "" || apiKey == "" {
		return "", domain.ErrInvalidAIConfig
	}

	endpoint, err := chatCompletionEndpoint(provider)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(openAIChatCompletionRequest{
		Model:    model,
		Messages: openAICompatibleMessages(request.Messages),
	})
	if err != nil {
		return "", err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+apiKey)

	response, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("chat completion failed: %s", aiErrorMessage(responseBody, response.StatusCode))
	}

	var completion openAIChatCompletionResponse
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return "", err
	}
	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("chat completion returned no choices")
	}

	content := strings.TrimSpace(completion.Choices[0].Message.Content)
	if content == "" {
		return "", fmt.Errorf("chat completion returned empty content")
	}

	return content, nil
}

func chatCompletionEndpoint(provider string) (string, error) {
	switch provider {
	case "openai":
		return "https://api.openai.com/v1/chat/completions", nil
	case "deepseek":
		return "https://api.deepseek.com/chat/completions", nil
	default:
		return "", domain.ErrUnsupportedAIProvider
	}
}

func openAICompatibleMessages(messages []domain.ChatMessage) []openAIChatMessage {
	compatibleMessages := make([]openAIChatMessage, 0, len(messages))
	for _, message := range messages {
		role := strings.TrimSpace(message.Role)
		content := strings.TrimSpace(message.Content)
		if role == "" || content == "" {
			continue
		}

		compatibleMessages = append(compatibleMessages, openAIChatMessage{
			Role:    role,
			Content: content,
		})
	}

	return compatibleMessages
}

func aiErrorMessage(responseBody []byte, statusCode int) string {
	var response openAIErrorResponse
	if err := json.Unmarshal(responseBody, &response); err == nil && response.Error.Message != "" {
		return response.Error.Message
	}

	return fmt.Sprintf("status %d", statusCode)
}

type openAIChatCompletionRequest struct {
	Model    string              `json:"model"`
	Messages []openAIChatMessage `json:"messages"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatCompletionResponse struct {
	Choices []struct {
		Message openAIChatMessage `json:"message"`
	} `json:"choices"`
}

type openAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}
