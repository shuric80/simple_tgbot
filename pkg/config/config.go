package config

import (
	"log"
	"os"
)

// AuthType тип аутентификации
type AuthType string

const (
	AuthTypeStaticToken    AuthType = "static_token"
	AuthTypeServiceAccount AuthType = "service_account"
	AuthTypeOAuth          AuthType = "oauth"
)

// Config структура конфигурации приложения
type Config struct {
	TelegramBotToken      string
	YandexCatalogID       string
	YandexCloudToken      string   // Статический токен (опционально)
	AuthType              AuthType // Тип аутентификации
	ServiceAccountKeyFile string   // Путь к файлу ключа сервисного аккаунта
	ServiceAccountKeyJSON string   // JSON ключа сервисного аккаунта
	YandexOAuthToken      string   // OAuth токен Яндекс Паспорта
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	cfg := &Config{
		TelegramBotToken: getEnvRequired("TELEGRAM_BOT_TOKEN"),
		YandexCatalogID:  getEnvRequired("YANDEX_CLOUD_CATALOG_ID"),
		YandexCloudToken: os.Getenv("YANDEX_CLOUD_TOKEN"),
		AuthType:         AuthType(getEnvWithDefault("YANDEX_AUTH_TYPE", "static_token")),
	}

	// Проверяем тип аутентификации и загружаем соответствующие параметры
	switch cfg.AuthType {
	case AuthTypeStaticToken:
		if cfg.YandexCloudToken == "" {
			log.Fatal("YANDEX_CLOUD_TOKEN is required for static_token auth type")
		}
	case AuthTypeServiceAccount:
		cfg.ServiceAccountKeyFile = os.Getenv("YANDEX_SERVICE_ACCOUNT_KEY_FILE")
		cfg.ServiceAccountKeyJSON = os.Getenv("YANDEX_SERVICE_ACCOUNT_KEY_JSON")
		if cfg.ServiceAccountKeyFile == "" && cfg.ServiceAccountKeyJSON == "" {
			log.Fatal("YANDEX_SERVICE_ACCOUNT_KEY_FILE or YANDEX_SERVICE_ACCOUNT_KEY_JSON is required for service_account auth type")
		}
	case AuthTypeOAuth:
		cfg.YandexOAuthToken = getEnvRequired("YANDEX_OAUTH_TOKEN")
	default:
		log.Fatalf("Unknown auth type: %s. Supported types: static_token, service_account, oauth", cfg.AuthType)
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

// getEnvWithDefault получает переменную окружения с значением по умолчанию
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
