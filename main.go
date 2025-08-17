package main

import (
	"log"
	"tgbot/pkg/bot"
	"tgbot/pkg/config"
	"tgbot/pkg/gpt"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем клиент для Yandex GPT
	gptClient := gpt.NewClient(cfg.YandexCatalogID, cfg.YandexCloudToken)

	// Создаем экземпляр бота
	telegramBot, err := bot.New(cfg.TelegramBotToken, gptClient)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Запускаем бота
	if err := telegramBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
