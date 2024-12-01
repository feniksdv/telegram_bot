package services

import (
	"bot/internal/app/helper"
	"bot/internal/app/message"
	"bot/internal/database/category_incomes"
	"bot/internal/database/customers"
	"bot/internal/database/session_bots"
	"fmt"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Incomes(update tgbotapi.Update) (string, *tgbotapi.ReplyKeyboardMarkup) {
	text := update.Message.Text

	text = strings.Replace(text, "+", "", 1)
	text = strings.Replace(text, ",", ".", -1)

	// парсим строку
	exp, err := govaluate.NewEvaluableExpression(text)
	if err != nil {
		//fmt.Println("Ошибка парсинга:", err)
		return message.ErrorIncomingData(update), nil
	}

	// Вычисление результата
	result, err := exp.Evaluate(nil)
	if err != nil {
		//fmt.Println("Ошибка вычисления:", err)
		return message.ErrorIncomingData(update), nil
	}

	// Деление на ноль
	if helper.CheckType(result) {
		return message.ErrorIncomingData(update), nil
	}

	resultStr := fmt.Sprintf("%v", result)

	categoryIncomes := createSessionBotStep1(update, resultStr)

	return resultStr, categoryIncomes
}

// получаем пользователя и создаем сессию для первого шага
func createSessionBotStep1(update tgbotapi.Update, resultStr string) *tgbotapi.ReplyKeyboardMarkup {
	// получить пользователя по telegram_id
	customer := customers.GetByCustomerTelegramId(update.Message.From.ID)

	var step uint8 = 1
	now := time.Now()
	typeValue := "incomes"
	userId := customer.UserID

	sessionBots := helper.SessionBots{
		UserTelegramId: &update.Message.From.ID,
		UserId:         userId,
		Money:          &resultStr,
		MoneyMessageId: &update.Message.MessageID,
		Type:           &typeValue,
		UnitId:         &step,
		Step:           &step,
		CreatedAt:      &now,
		UpdatedAt:      &now,
		DeletedAt:      &now,
	}
	session_bots.Create(sessionBots)

	categoryIncomes, err := category_incomes.GetCategoriesIncomesByUserId(userId)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return createKeyboard(categoryIncomes)
}

func createKeyboard(categoryIncomes []helper.CategoryIncomes) *tgbotapi.ReplyKeyboardMarkup {
	if len(categoryIncomes) == 0 {
		return nil
	}

	var rows [][]tgbotapi.KeyboardButton
	var currentRow []tgbotapi.KeyboardButton

	for _, income := range categoryIncomes {
		if income.Id != nil && income.Name != nil {
			button := tgbotapi.NewKeyboardButton(*income.Name)
			currentRow = append(currentRow, button)

			// Если в текущей строке уже 3 кнопки, добавляем строку в rows и сбрасываем currentRow
			if len(currentRow) == 3 {
				rows = append(rows, tgbotapi.NewKeyboardButtonRow(currentRow...))
				currentRow = nil // Сбрасываем текущую строку
			}
		}
	}

	// Если остались кнопки в currentRow, добавляем их как последнюю строку
	if len(currentRow) > 0 {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(currentRow...))
	}

	// Создаем ReplyKeyboardMarkup с помощью конструктора и берем адрес его значения
	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	return &keyboard
}
