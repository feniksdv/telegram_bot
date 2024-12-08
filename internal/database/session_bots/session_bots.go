package session_bots

import (
	"bot/internal/app/helper"
	"bot/internal/database"
	"bot/internal/database/expenses"
	"bot/internal/database/incomes"
	"database/sql"
	"fmt"
	"time"
)

func GetByUserTelegramId(userTelegramID int64) *helper.SessionBots {
	db := database.Connect()
	defer db.Close()

	query := `
        SELECT 
            sb.user_telegram_id, 
            sb.user_id, 
			money,
            sb.step, 
            sb.created_at, 
            sb.deleted_at
        FROM session_bots AS sb 
        LEFT JOIN users AS u ON u.id = sb.user_id 
        WHERE sb.deleted_at IS NULL and user_telegram_id = ?
    `

	results, err := db.Query(query, userTelegramID)
	if err != nil {
		panic(err.Error()) // В реальном приложении используйте более надёжную обработку ошибок
	}
	defer results.Close()

	var sessionBots helper.SessionBots
	var createdAtStr, deletedAtStr sql.NullString

	for results.Next() {
		err = results.Scan(
			&sessionBots.UserTelegramId,
			&sessionBots.UserId,
			&sessionBots.Money,
			&sessionBots.Step,
			&createdAtStr,
			&deletedAtStr,
		)
		if err != nil {
			panic(err.Error())
		}

		// Преобразуем created_at из строки в time.Time
		if createdAtStr.Valid {
			createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr.String)
			if err != nil {
				panic(fmt.Sprintf("Ошибка парсинга created_at: %v", err))
			}
			sessionBots.CreatedAt = &createdAt
		}

		// Преобразуем deleted_at из строки в time.Time (если нужно)
		if deletedAtStr.Valid {
			deletedAt, err := time.Parse("2006-01-02 15:04:05", deletedAtStr.String)
			if err != nil {
				panic(fmt.Sprintf("Ошибка парсинга deleted_at: %v", err))
			}
			sessionBots.DeletedAt = &deletedAt
		}
	}

	if sessionBots.UserId != nil {
		return &sessionBots
	}
	return nil
}

func GetFirst() *helper.SessionBots {
	db := database.Connect()
	defer db.Close()
	results, err := db.Query("SELECT user_telegram_id, user_id FROM session_bots")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var user helper.SessionBots

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&user.UserTelegramId, &user.UserId)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	// and then print out the tag's Name attribute
	return &user
}

func Create(sessionBots helper.SessionBots) {
	db := database.Connect()
	defer db.Close()

	results, err := db.Prepare("INSERT INTO session_bots(user_telegram_id,user_id,money,money_message_id,type,unit_id,step,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = results.Exec(
		sessionBots.UserTelegramId,
		sessionBots.UserId,
		sessionBots.Money,
		sessionBots.MoneyMessageId,
		sessionBots.Type,
		sessionBots.UnitId,
		sessionBots.Step,
		sessionBots.CreatedAt,
		sessionBots.UpdatedAt)
	if err != nil {
		panic(err.Error())
	}
}

func UpdateStep2(messageId int, categoryName string, categoryId uint64, telegramUserId int64, userId *int32) error {
	db := database.Connect()
	defer db.Close()

	query := `
			UPDATE session_bots
			SET category_id = ?,
				category_name = ?,
				category_message_id = ?,
				step = 2
			WHERE user_telegram_id = ? and user_id = ? and deleted_at IS NULL and step = 1
		`

	_, err := db.Exec(query, categoryId, categoryName, messageId, telegramUserId, *userId)
	if err != nil {
		return fmt.Errorf("failed to update session_bots: %w", err)
	}

	return nil
}

func UpdateStep3(messageId int, categoryName string, categoryId uint64, telegramUserId int64, userId *int32) error {
	db := database.Connect()
	defer db.Close()

	query := `
			UPDATE session_bots
			SET category_source_id = ?,
				category_source_name = ?,
				source_message_id = ?,
				step = 3
			WHERE user_telegram_id = ? and user_id = ? and deleted_at IS NULL and step = 2
		`

	_, err := db.Exec(query, categoryId, categoryName, messageId, telegramUserId, *userId)
	if err != nil {
		return fmt.Errorf("failed to update session_bots: %w", err)
	}

	return nil
}

func UpdateStep4(messageId int, description string, telegramUserId int64, userId *int32) error {
	db := database.Connect()
	defer db.Close()

	query := `
			UPDATE session_bots
			SET description = ?,
				description_message_id = ?,
				step = 4
			WHERE user_telegram_id = ? and user_id = ? and deleted_at IS NULL and step = 3
		`

	_, err := db.Exec(query, description, messageId, telegramUserId, userId)
	if err != nil {
		return fmt.Errorf("failed to update session_bots: %w", err)
	}

	return nil
}

func UpdateStep5(telegramUserId int64, userId *int32) error {
	db := database.Connect()
	defer db.Close()

	q := `
			SELECT user_telegram_id,
			user_id,
			money,
			money_message_id,
			category_id,
			category_name,
			category_message_id,
			type,
			unit_id,
			description,
			description_message_id,
			category_source_id,
			category_source_name,
			source_message_id 
			FROM session_bots
			WHERE user_telegram_id = ? and user_id = ? and deleted_at IS NULL and step = 4 
	`
	results, e := db.Query(q, telegramUserId, userId)
	if e != nil {
		panic(e.Error())
	}
	defer results.Close() // Закрываем результаты после использования

	var sessionBot helper.SessionBots
	if results.Next() {
		if err := results.Scan(
			&sessionBot.UserTelegramId,
			&sessionBot.UserId,
			&sessionBot.Money,
			&sessionBot.MoneyMessageId,
			&sessionBot.CategoryId,
			&sessionBot.CategoryName,
			&sessionBot.CategoryMessageId,
			&sessionBot.Type,
			&sessionBot.UnitId,
			&sessionBot.Description,
			&sessionBot.DescriptionMessageId,
			&sessionBot.CategorySourceId,
			&sessionBot.CategorySourceName,
			&sessionBot.SourceMessageId); err != nil {
			return fmt.Errorf("failed to scan type: %w", err)
		}
	} else {
		return fmt.Errorf("no rows found")
	}

	query := `
			UPDATE session_bots
			SET deleted_at = CURRENT_TIMESTAMP
			WHERE user_telegram_id = ? and user_id = ? and deleted_at IS NULL and step = 4
		`

	_, err := db.Exec(query, telegramUserId, userId)
	if err != nil {
		panic(err.Error())
	}

	if *sessionBot.Type == "incomes" {
		incomes.Create(sessionBot)
	} else {
		expenses.Create(sessionBot)
	}
	return nil
}

func Insert() {

}

func Delete() {

}
