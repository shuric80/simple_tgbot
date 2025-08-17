# Аутентификация с помощью IAM токенов Yandex Cloud

В проекте добавлена поддержка динамического получения IAM токенов для работы с Yandex Cloud API. Это позволяет использовать более безопасные методы аутентификации вместо статических токенов.

## Поддерживаемые методы аутентификации

### 1. Статический токен (по умолчанию)
Используется для обратной совместимости. Токен передается напрямую без обновления.

**Переменные окружения:**
```bash
YANDEX_AUTH_TYPE=static_token
YANDEX_CLOUD_TOKEN=<your_static_token>
```

### 2. Service Account Key (рекомендуется)
Использует ключ сервисного аккаунта для автоматического получения и обновления IAM токенов.

**Переменные окружения:**
```bash
YANDEX_AUTH_TYPE=service_account

# Опция 1: Путь к файлу с ключом
YANDEX_SERVICE_ACCOUNT_KEY_FILE=/path/to/service-account-key.json

# Опция 2: JSON ключа в переменной окружения
YANDEX_SERVICE_ACCOUNT_KEY_JSON='{"id":"...","service_account_id":"...","private_key":"..."}'
```

### 3. OAuth Token
Использует OAuth токен Яндекс Паспорта для получения IAM токенов.

**Переменные окружения:**
```bash
YANDEX_AUTH_TYPE=oauth
YANDEX_OAUTH_TOKEN=<your_oauth_token>
```

## Структура ключа сервисного аккаунта

Файл ключа сервисного аккаунта должен содержать:
```json
{
  "id": "ключевой_идентификатор",
  "service_account_id": "идентификатор_сервисного_аккаунта",
  "created_at": "2023-01-01T00:00:00Z",
  "key_algorithm": "RSA_2048",
  "public_key": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
}
```

## Как получить ключ сервисного аккаунта

1. Создайте сервисный аккаунт в Yandex Cloud Console
2. Назначьте необходимые роли (например, `ai.languageModels.user`)
3. Создайте ключ для сервисного аккаунта:
   ```bash
   yc iam key create --service-account-name <service-account-name> --output key.json
   ```

## Преимущества IAM токенов

- **Безопасность**: Токены автоматически обновляются и имеют ограниченное время жизни
- **Контроль доступа**: Можно точно настроить права сервисного аккаунта
- **Аудит**: Все операции логируются в Yandex Cloud
- **Ротация**: Ключи можно безопасно ротировать без изменения кода

## Переход с статических токенов

Для перехода на IAM токены:

1. Создайте сервисный аккаунт и ключ
2. Измените переменную `YANDEX_AUTH_TYPE=service_account`
3. Добавьте `YANDEX_SERVICE_ACCOUNT_KEY_FILE` или `YANDEX_SERVICE_ACCOUNT_KEY_JSON`
4. Удалите `YANDEX_CLOUD_TOKEN` (опционально)

## Обработка ошибок

Система автоматически:
- Обновляет токены за 5 минут до истечения
- Повторно получает токены при ошибках аутентификации
- Логирует все операции с токенами

## Пример docker-compose.yml

```yaml
version: '3.8'
services:
  tgbot:
    build: .
    environment:
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
      - YANDEX_CLOUD_CATALOG_ID=${YANDEX_CLOUD_CATALOG_ID}
      - YANDEX_AUTH_TYPE=service_account
      - YANDEX_SERVICE_ACCOUNT_KEY_FILE=/app/service-account-key.json
    volumes:
      - ./service-account-key.json:/app/service-account-key.json:ro
```

## Troubleshooting

### Ошибки аутентификации
- Проверьте права сервисного аккаунта
- Убедитесь, что ключ не поврежден
- Проверьте синхронизацию времени на сервере

### Ошибки JWT
- Проверьте формат приватного ключа (должен быть PKCS#8)
- Убедитесь, что поля `id` и `service_account_id` корректны

### Сетевые ошибки
- Проверьте доступность `iam.api.cloud.yandex.net`
- Убедитесь в корректности сетевых настроек
