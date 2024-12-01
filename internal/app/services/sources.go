package services

import (
	"bot/internal/app/helper"
	"bot/internal/database/category_sources"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Sources(userId *int32) tgbotapi.ReplyKeyboardMarkup {
	categorySources, err := category_sources.GetAllCategorySourcesByUserId(userId)
	if err != nil {
		fmt.Println("Error:", err)
		return tgbotapi.ReplyKeyboardMarkup{}
	}

	return createKeyboardSources(categorySources)
}

// подготовим клавиатуру
func createKeyboardSources(categorySources []helper.CategorySource) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	var currentRow []tgbotapi.KeyboardButton

	for _, source := range categorySources {
		if source.ID != nil && source.Name != nil {
			button := tgbotapi.NewKeyboardButton(*source.Name)
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

	// Возвращаем клавиатуру
	return tgbotapi.NewReplyKeyboard(rows...)
}
