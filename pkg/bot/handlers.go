package bot

import (
	"log"
	"tgbot/pkg/gpt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot структура для работы с Telegram ботом
type Bot struct {
	api       *tgbotapi.BotAPI
	gptClient *gpt.Client
}

// New создает новый экземпляр бота
func New(token string, gptClient *gpt.Client) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	// Включаем отладку (опционально)
	api.Debug = true

	return &Bot{
		api:       api,
		gptClient: gptClient,
	}, nil
}

// Start запускает бота
func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	// Настраиваем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	// Обрабатываем входящие сообщения
	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}

	return nil
}

// handleMessage обрабатывает входящие сообщения
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	// Отправляем сообщение в Yandex GPT
	gptResponse, err := b.gptClient.SendMessage(message.Text)
	if err != nil {
		log.Printf("Error getting GPT response: %v", err)
		// В случае ошибки отправляем эхо
		b.sendErrorMessage(message, err)
		return
	}

	// Создаем ответное сообщение с ответом от GPT
	b.sendResponse(message, gptResponse)
}

// sendErrorMessage отправляет сообщение об ошибке
func (b *Bot) sendErrorMessage(message *tgbotapi.Message, err error) {
	errorText := "Извините, произошла ошибка при обработке вашего сообщения. Ваше сообщение: " + message.Text
	msg := tgbotapi.NewMessage(message.Chat.ID, errorText)
	msg.ReplyToMessageID = message.MessageID

	if _, sendErr := b.api.Send(msg); sendErr != nil {
		log.Printf("Error sending error message: %v", sendErr)
	}
}

// sendResponse отправляет ответ пользователю
func (b *Bot) sendResponse(message *tgbotapi.Message, response string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ReplyToMessageID = message.MessageID

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
