package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"omni-pixel/domain"
)

type OpenAIProvider struct {
	BaseURL string
	APIKey  string
}

func NewOpenAIProvider(baseURL, apiKey string) *OpenAIProvider {
	return &OpenAIProvider{BaseURL: baseURL, APIKey: apiKey}
}

type chatRequest struct {
	Model    string               `json:"model"`
	Messages []domain.AIChatMessage `json:"messages"`
	Stream   bool                 `json:"stream"`
}

type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

func (p *OpenAIProvider) ChatStream(messages []domain.AIChatMessage, modelID string) (<-chan domain.AIStreamChunk, error) {
	body := chatRequest{
		Model:    modelID,
		Messages: messages,
		Stream:   true,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", p.BaseURL+"/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan domain.AIStreamChunk)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					ch <- domain.AIStreamChunk{Done: true}
				}
				return
			}

			line = strings.TrimSpace(line)
			if line == "" || !strings.HasPrefix(line, "data:") {
				continue
			}

			data := strings.TrimPrefix(line, "data:")
			data = strings.TrimSpace(data)
			if data == "[DONE]" {
				return
			}

			var chunk streamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				ch <- domain.AIStreamChunk{Token: chunk.Choices[0].Delta.Content}
			}

			if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason != nil {
				return
			}
		}
	}()

	return ch, nil
}
