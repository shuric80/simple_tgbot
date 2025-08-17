package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tgbot/pkg/iam"
)

// Request структура запроса к Yandex GPT
type Request struct {
	ModelURI          string            `json:"modelUri"`
	CompletionOptions CompletionOptions `json:"completionOptions"`
	Messages          []Message         `json:"messages"`
}

// CompletionOptions опции для генерации
type CompletionOptions struct {
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
}

// Message сообщение для GPT
type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

// Response структура ответа от Yandex GPT
type Response struct {
	Result Result `json:"result"`
}

// Result результат ответа
type Result struct {
	Alternatives []Alternative `json:"alternatives"`
}

// Alternative альтернативный ответ
type Alternative struct {
	Message Message `json:"message"`
	Status  string  `json:"status"`
}

// Client клиент для работы с Yandex GPT API
type Client struct {
	catalogID    string
	token        string
	tokenManager *iam.TokenManager
	client       *http.Client
}

// NewClient создает новый клиент для Yandex GPT со статическим токеном
func NewClient(catalogID, token string) *Client {
	return &Client{
		catalogID: catalogID,
		token:     token,
		client:    &http.Client{},
	}
}

// NewClientWithTokenManager создает новый клиент для Yandex GPT с менеджером токенов
func NewClientWithTokenManager(catalogID string, tokenManager *iam.TokenManager) *Client {
	return &Client{
		catalogID:    catalogID,
		tokenManager: tokenManager,
		client:       &http.Client{},
	}
}

// SendMessage отправляет сообщение в Yandex GPT и возвращает ответ
func (c *Client) SendMessage(userMessage string) (string, error) {
	modelURI := fmt.Sprintf("gpt://%s/yandexgpt-lite", c.catalogID)

	request := Request{
		ModelURI: modelURI,
		CompletionOptions: CompletionOptions{
			Stream:      false,
			Temperature: 0.6,
			MaxTokens:   2000,
		},
		Messages: []Message{
			{
				Role: "user",
				Text: userMessage,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Получаем токен
	token, err := c.getToken()
	if err != nil {
		return "", fmt.Errorf("error getting token: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(response.Result.Alternatives) == 0 {
		return "", fmt.Errorf("no alternatives in response")
	}

	return response.Result.Alternatives[0].Message.Text, nil
}

// getToken возвращает токен для аутентификации
func (c *Client) getToken() (string, error) {
	// Если есть менеджер токенов, используем его
	if c.tokenManager != nil {
		return c.tokenManager.GetToken()
	}

	// Иначе используем статический токен
	if c.token == "" {
		return "", fmt.Errorf("no token available")
	}

	return c.token, nil
}
