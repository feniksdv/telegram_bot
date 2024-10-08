package services

import (
	"bot/internal/app/helper"
	"bot/internal/app/message"
	"fmt"
	"github.com/Knetic/govaluate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func Incomes(update tgbotapi.Update) string {
	text := update.Message.Text

	fmt.Println("Incomes: " + text)
	text = strings.Replace(text, "+", "", 1)
	text = strings.Replace(text, ",", ".", -1)

	// парсим строку
	exp, err := govaluate.NewEvaluableExpression(text)
	if err != nil {
		//fmt.Println("Ошибка парсинга:", err)
		return message.ErrorIncomingData(update)
	}

	// Вычисление результата
	result, err := exp.Evaluate(nil)
	if err != nil {
		//fmt.Println("Ошибка вычисления:", err)
		return message.ErrorIncomingData(update)
	}

	// Деление на ноль
	if helper.CheckType(result) {
		return message.ErrorIncomingData(update)
	}

	resultStr := fmt.Sprintf("%v", result)

	return resultStr
}
