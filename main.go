package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Получаем токен бота из переменной окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

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

			// Создаем ответное сообщение (эхо)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// Отвечаем на то же сообщение
			msg.ReplyToMessageID = update.Message.MessageID

			// Отправляем сообщение
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
