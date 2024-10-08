package helper

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"regexp"
)

// определить начинается ли строка с цифры или (
func CheckIfStartsWithDigit(update tgbotapi.Update) bool {
	matched, _ := regexp.MatchString("^[0-9(]", update.Message.Text)
	return matched
}

// проверить что выражение не равняется бесконечности
func CheckType(x interface{}) bool {
	fx := x.(float64)
	result := fx == math.Inf(1) || fx == -math.Inf(1)
	return result
}
