package main

import (
	"log"
	"tgbot/pkg/bot"
	"tgbot/pkg/config"
	"tgbot/pkg/gpt"
	"tgbot/pkg/iam"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем клиент для Yandex GPT в зависимости от типа аутентификации
	var gptClient *gpt.Client

	switch cfg.AuthType {
	case config.AuthTypeStaticToken:
		gptClient = gpt.NewClient(cfg.YandexCatalogID, cfg.YandexCloudToken)
		log.Println("Using static token authentication")

	case config.AuthTypeServiceAccount:
		var serviceAccountKey *iam.ServiceAccountKey
		var err error

		if cfg.ServiceAccountKeyFile != "" {
			serviceAccountKey, err = iam.LoadServiceAccountKeyFromFile(cfg.ServiceAccountKeyFile)
			if err != nil {
				log.Fatalf("Failed to load service account key from file: %v", err)
			}
			log.Printf("Loaded service account key from file: %s", cfg.ServiceAccountKeyFile)
		} else {
			serviceAccountKey, err = iam.LoadServiceAccountKeyFromJSON(cfg.ServiceAccountKeyJSON)
			if err != nil {
				log.Fatalf("Failed to load service account key from JSON: %v", err)
			}
			log.Println("Loaded service account key from environment variable")
		}

		tokenManager := iam.NewTokenManagerWithServiceAccount(serviceAccountKey)
		gptClient = gpt.NewClientWithTokenManager(cfg.YandexCatalogID, tokenManager)
		log.Println("Using service account authentication")

	case config.AuthTypeOAuth:
		tokenManager := iam.NewTokenManagerWithOAuth(cfg.YandexOAuthToken)
		gptClient = gpt.NewClientWithTokenManager(cfg.YandexCatalogID, tokenManager)
		log.Println("Using OAuth token authentication")

	default:
		log.Fatalf("Unknown authentication type: %s", cfg.AuthType)
	}

	// Создаем экземпляр бота
	telegramBot, err := bot.New(cfg.TelegramBotToken, gptClient)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Запускаем бота
	log.Println("Starting Telegram bot...")
	if err := telegramBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
