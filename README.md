# Telegram Bot на Go с Yandex GPT

Intelligent Telegram бот с интеграцией Yandex GPT API, написанный на Go и запускаемый через Docker Compose.

## Функционал

- **Интеграция с Yandex GPT** - каждое сообщение обрабатывается через Yandex GPT API
- **Умные ответы** - бот отвечает на основе анализа сообщения с помощью ИИ
- **Fallback в эхо-режим** - в случае ошибки API возвращается к простому эхо
- Логирование всех входящих сообщений и ответов API
- Готов к запуску в Docker контейнере

## Настройка

1. **Создайте бота в Telegram:**
   - Напишите [@BotFather](https://t.me/BotFather) в Telegram
   - Создайте нового бота командой `/newbot`
   - Получите токен бота

2. **Настройте Yandex Cloud:**
   - Создайте аккаунт в [Yandex Cloud](https://cloud.yandex.ru/)
   - Получите Catalog ID из консоли Yandex Cloud
   - Создайте API токен для доступа к Yandex GPT

3. **Создайте файл `.env` в корне проекта:**
   ```bash
   TELEGRAM_BOT_TOKEN=ваш_токен_от_BotFather
   YANDEX_CLOUD_CATALOG_ID=ваш_catalog_id
   YANDEX_CLOUD_TOKEN=ваш_yandex_api_token
   ```

4. **Запустите бота:**
   ```bash
   docker-compose up --build
   ```

## Как работает

1. **Пользователь отправляет сообщение** боту в Telegram
2. **Бот получает сообщение** и логирует его
3. **Сообщение отправляется** в Yandex GPT API для обработки
4. **Yandex GPT анализирует** текст и генерирует умный ответ
5. **Бот отправляет** ответ от GPT пользователю
6. **В случае ошибки** API бот работает в эхо-режиме

## Переменные окружения

- `TELEGRAM_BOT_TOKEN` - токен бота от BotFather (обязательно)
- `YANDEX_CLOUD_CATALOG_ID` - ID каталога Yandex Cloud (обязательно)  
- `YANDEX_CLOUD_TOKEN` - API токен Yandex Cloud (обязательно)

## Структура проекта

```
tgbot/
├── main.go              # Основной код бота
├── go.mod               # Go модули
├── Dockerfile           # Docker образ
├── docker-compose.yml   # Docker Compose конфигурация
├── .env                 # Переменные окружения (создать вручную)
└── README.md           # Этот файл
```

## Команды для разработки

### Запуск локально (без Docker)
```bash
go mod tidy
export TELEGRAM_BOT_TOKEN="ваш_токен"
go run main.go
```

### Запуск в Docker
```bash
docker-compose up --build
```

### Остановка
```bash
docker-compose down
```

### Просмотр логов
```bash
docker-compose logs -f tgbot
```

## Переменные окружения

- `TELEGRAM_BOT_TOKEN` - токен бота от BotFather (обязательно)

## Следующие шаги

Этот базовый бот можно расширить:
- Добавить команды (/start, /help и т.д.)
- Подключить базу данных
- Добавить middleware для логирования
- Реализовать более сложную логику обработки сообщений
