package commands

import (
	"bot/internal/app/helper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Default - дефолтная команда
func Default(update tgbotapi.Update) string {
	name := helper.GetContactName(update)
	return "Привет, " + name + "\nДанные не сохранены, вы ввели не существующею валюту."
}
