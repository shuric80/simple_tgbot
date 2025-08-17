package config

import (
	"log"
	"os"
)

// Config структура конфигурации приложения
type Config struct {
	TelegramBotToken   string
	YandexCatalogID    string
	YandexCloudToken   string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	cfg := &Config{
		TelegramBotToken: getEnvRequired("TELEGRAM_BOT_TOKEN"),
		YandexCatalogID:  getEnvRequired("YANDEX_CLOUD_CATALOG_ID"),
		YandexCloudToken: getEnvRequired("YANDEX_CLOUD_TOKEN"),
	}

	return cfg
}

// getEnvRequired получает обязательную переменную окружения
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable is required", key)
	}
	return value
}
