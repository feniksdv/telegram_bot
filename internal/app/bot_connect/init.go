package bot_connect

import (
	"bot/internal/app/helper"
	"bot/internal/app/message"
	"bot/internal/app/services"
	"bot/internal/database"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

// Init инициализация бота
func Init() {

	updates, bot := connect()

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			// открываем ссесию записи данных
			// session_bot - поиск по user_telegram_id и deleted_at = null если нет, то запускаем step 1
			// иначе вернуть step и запустить логику ниже

			query := database.GetByUserTelegramId()
			if query == nil {
				fmt.Println("null")
			} else {
				fmt.Println("user ID = " + *query)
			}
			// step 1
			// создаем запиьс в session_bot
			// из customers и users получаем нужные поля
			// сохранить расходи или доход id операции
			// возвращаем step и запускаем логику ниже
			switch {
			// если в начале стоит плюс
			case strings.HasPrefix(update.Message.Text, "+"):
				result := services.Incomes(update)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, result)
			// если в начале число
			case helper.CheckIfStartsWithDigit(update):
				result := services.Expenses(update)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, result)
			default:
				errorIncomingData := message.ErrorIncomingData(update)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, errorIncomingData)
			}
			// первый шаг вернуть должен набор кнопок с категориями - в зависимости расходы доходы

			// step 2
			// получает значение кнопки - сохраняет где надо и возвращает сообщение введите описание
			// возвращает набор откуда производились траты

			// step 3
			// получаем id трат и сохроняем
			// закрываем ссесию

			msg.ReplyMarkup = numericKeyboard
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
