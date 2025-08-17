package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// YandexGPTRequest структура запроса к Yandex GPT
type YandexGPTRequest struct {
	ModelURI          string      `json:"modelUri"`
	CompletionOptions CompletionOptions `json:"completionOptions"`
	Messages          []GPTMessage `json:"messages"`
}

// CompletionOptions опции для генерации
type CompletionOptions struct {
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
}

// GPTMessage сообщение для GPT
type GPTMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

// YandexGPTResponse структура ответа от Yandex GPT
type YandexGPTResponse struct {
	Result Result `json:"result"`
}

type Result struct {
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Message GPTMessage `json:"message"`
	Status  string     `json:"status"`
}

// YandexGPTClient клиент для работы с Yandex GPT API
type YandexGPTClient struct {
	catalogID string
	token     string
	client    *http.Client
}

// NewYandexGPTClient создает новый клиент для Yandex GPT
func NewYandexGPTClient(catalogID, token string) *YandexGPTClient {
	return &YandexGPTClient{
		catalogID: catalogID,
		token:     token,
		client:    &http.Client{},
	}
}

// SendMessage отправляет сообщение в Yandex GPT и возвращает ответ
func (c *YandexGPTClient) SendMessage(userMessage string) (string, error) {
	modelURI := fmt.Sprintf("gpt://%s/yandexgpt-lite", c.catalogID)
	
	request := YandexGPTRequest{
		ModelURI: modelURI,
		CompletionOptions: CompletionOptions{
			Stream:      false,
			Temperature: 0.6,
			MaxTokens:   2000,
		},
		Messages: []GPTMessage{
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
	req.Header.Set("Authorization", "Bearer "+c.token)

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

	var response YandexGPTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(response.Result.Alternatives) == 0 {
		return "", fmt.Errorf("no alternatives in response")
	}

	return response.Result.Alternatives[0].Message.Text, nil
}

func main() {
	// Получаем токен бота из переменной окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// Получаем параметры Yandex Cloud
	catalogID := os.Getenv("YANDEX_CLOUD_CATALOG_ID")
	if catalogID == "" {
		log.Fatal("YANDEX_CLOUD_CATALOG_ID environment variable is required")
	}

	yandexToken := os.Getenv("YANDEX_CLOUD_TOKEN")
	if yandexToken == "" {
		log.Fatal("YANDEX_CLOUD_TOKEN environment variable is required")
	}

	// Создаем клиент для Yandex GPT
	gptClient := NewYandexGPTClient(catalogID, yandexToken)

	// Создаем новый экземпляр бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Включаем отладку (опционально)
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Настраиваем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обрабатываем входящие сообщения
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Отправляем сообщение в Yandex GPT
			gptResponse, err := gptClient.SendMessage(update.Message.Text)
			if err != nil {
				log.Printf("Error getting GPT response: %v", err)
				// В случае ошибки отправляем эхо
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, произошла ошибка при обработке вашего сообщения. Ваше сообщение: "+update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				continue
			}

			// Создаем ответное сообщение с ответом от GPT
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, gptResponse)

			// Отвечаем на то же сообщение
			msg.ReplyToMessageID = update.Message.MessageID

			// Отправляем сообщение
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
