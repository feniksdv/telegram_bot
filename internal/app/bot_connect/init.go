package bot_connect

import (
	"bot/internal/app/commands"
	"bot/internal/app/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

// Init инициализация бота
func Init() {

	updates, bot := connect()

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Text {
			case "/help":
				help := commands.Help(update)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, help)
			default:
				defaultMsg := commands.Default(update)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, defaultMsg)
			}

			switch {
			// если в начале стоит плюс
			case strings.HasPrefix(update.Message.Text, "+"):
				result := services.Calc(update.Message.Text)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, result)
				// определить что ввел пользователь

				// если это сложение чисел
				// если это расход
			}

			bot.Send(msg)
		}
	}
}

// Подключение к боту Telegram
func connect() (tgbotapi.UpdatesChannel, *tgbotapi.BotAPI) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	botDebug := os.Getenv("BOT_DEBUG")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Инициализация бота завершилась ошибкой: ", err)
	}

	// Отладка бота
	bot.Debug, err = strconv.ParseBool(botDebug)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	config := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(config)

	return updates, bot
}
