package commands

import (
	"bot/internal/app/helper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Help(update tgbotapi.Update) string {
	name := helper.GetContactName(update)

	return "Привет " + name
}
