package bot_connect

import (
	"bot/internal/app/helper"
	"bot/internal/app/message"
	"bot/internal/app/services"
	"bot/internal/database/category_incomes"
	"bot/internal/database/category_sources"
	"bot/internal/database/customers"
	"bot/internal/database/session_bots"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// Init инициализация бота
func Init() {

	updates, bot := connect()

	for update := range updates {
		if update.Message != nil {

			var step uint8
			sessionBots := session_bots.GetByUserTelegramId(update.Message.From.ID)

			if sessionBots == nil {
				step = 1
			} else {
				step = *sessionBots.Step + 1
			}

			// осталось
			// сохранить в доходы или расходы после session_bots
			// реадктирование

			switch step {
			case 1:
				msg, keyboard := handleStep1(update)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			case 2:
				msg, keyboard := handleStep2(update)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			case 3:
				handleStep3(update, bot)
			case 4:
				handleStep4(update, bot)
			}
		}
	}
}

func connect() (tgbotapi.UpdatesChannel, *tgbotapi.BotAPI) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	botDebug := os.Getenv("BOT_DEBUG")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Инициализация бота завершилась ошибочно: ", err)
	}

	bot.Debug, _ = strconv.ParseBool(botDebug)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	config := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(config)

	return updates, bot
}

func handleStep1(update tgbotapi.Update) (tgbotapi.MessageConfig, *tgbotapi.ReplyKeyboardMarkup) {
	messageStr := "Выберите категорию доходов"

	switch {
	// если в начале стоит плюс
	case strings.HasPrefix(update.Message.Text, "+"):
		_, keyboard := services.Incomes(update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageStr)
		return msg, keyboard
	// если в начале число
	case helper.CheckIfStartsWithDigit(update):
		_, keyboard := services.Expenses(update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageStr)
		return msg, keyboard
	default:
		messageStr = message.ErrorIncomingData(update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageStr)
		return msg, nil
	}
}

func handleStep2(update tgbotapi.Update) (tgbotapi.MessageConfig, tgbotapi.ReplyKeyboardMarkup) {
	message := "Выберите форму оплаты:"
	customer := customers.GetByCustomerTelegramId(update.Message.From.ID)
	categories, err := category_incomes.GetCategoryIncomeByUserIdAndName(customer.UserID, update.Message.Text)
	if err != nil {
		message = fmt.Sprintf("Ошибка при получении категорий: %v, пришла строка %s", err, update.Message.Text)
	}
	session_bots.UpdateStep2(update.Message.MessageID, *categories.Name, *categories.Id, update.Message.From.ID, customer.UserID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	keyboard := services.Sources(customer.UserID)
	return msg, keyboard
}

func handleStep3(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	message := "Введите описание:"
	customer := customers.GetByCustomerTelegramId(update.Message.From.ID)
	categories, err := category_sources.GetCategorySourceByUserIdAndName(customer.UserID, update.Message.Text)
	if err != nil {
		message = fmt.Sprintf("Ошибка при получении категорий: %v, пришла строка %s", err, update.Message.Text)
	}
	session_bots.UpdateStep3(update.Message.MessageID, *categories.Name, *categories.ID, update.Message.From.ID, customer.UserID)
	closeKeyboard(bot, update.Message.Chat.ID, message)
}

func handleStep4(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	message := "Данные сохранены"
	customer := customers.GetByCustomerTelegramId(update.Message.From.ID)
	session_bots.UpdateStep4(update.Message.MessageID, update.Message.Text, update.Message.From.ID, customer.UserID)
	end(bot, update, update.Message.Chat.ID, message)
}

func closeKeyboard(bot *tgbotapi.BotAPI, chatID int64, message string) {
	replyMarkup := tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = replyMarkup

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func end(bot *tgbotapi.BotAPI, update tgbotapi.Update, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
	customer := customers.GetByCustomerTelegramId(update.Message.From.ID)
	session_bots.UpdateStep5(update.Message.From.ID, customer.UserID)
}
