package message

import (
	"bot/internal/app/helper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Выводит ссобщение о стандартном вводе
func ErrorIncomingData(update tgbotapi.Update) string {
	name := helper.GetContactName(update)
	return "Здравствуйте, " + name + "" +
		"\nДанные не сохранены, вы ввели не число или указали неверно формулу, пример разрешенных данных:" +
		"\n1. 1024," +
		"\n2. 1024.45," +
		"\n3. (1024.45+100+20)/(200+300)-1000" +
		"\nPS: Ошибка! В строке могут быть только цифры, скобки, и операторы +, -, *, /"
}
