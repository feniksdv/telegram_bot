package expenses

import (
	"bot/internal/app/helper"
	"bot/internal/database"
	"fmt"
)

func Create(sessionBot helper.SessionBots) {
	db := database.Connect()
	defer db.Close()

	query := `
		INSERT INTO expenses
		(	
			user_id,
			expense_message_id,
			category_expense_id,
			source_message_id,
			category_source_id,
			message_id,
			money,
			unit,
			unit_id,
			description_message_id,
			description,
			created_at,
			updated_at
		) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)
	`
	results, err := db.Prepare(query)
	if err != nil {
		fmt.Println(" Ошибка записи в expenses - ")
		panic(err.Error())
	}
	results.Exec(sessionBot.UserId, sessionBot.CategoryMessageId, sessionBot.CategoryId, sessionBot.SourceMessageId, sessionBot.CategorySourceId, sessionBot.MoneyMessageId, sessionBot.Money, "₽", sessionBot.UnitId, sessionBot.DescriptionMessageId, sessionBot.Description)
}
