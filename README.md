# Telegram Bot на Go

Простой эхо-бот для Telegram, написанный на Go и запускаемый через Docker Compose.

## Функционал

- Эхо-ответы на все входящие сообщения
- Логирование всех входящих сообщений
- Готов к запуску в Docker контейнере

## Настройка

1. **Создайте бота в Telegram:**
   - Напишите [@BotFather](https://t.me/BotFather) в Telegram
   - Создайте нового бота командой `/newbot`
   - Получите токен бота

2. **Создайте файл `.env` в корне проекта:**
   ```bash
   TELEGRAM_BOT_TOKEN=ваш_токен_от_BotFather
   ```

3. **Запустите бота:**
   ```bash
   docker-compose up --build
   ```

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
