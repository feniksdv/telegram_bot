package customers

import (
	"bot/internal/app/helper"
	"bot/internal/database"
)

func GetByCustomerTelegramId(userTelegramID int64) *helper.Customer {
	db := database.Connect()
	defer db.Close()

	query := `
        SELECT
            user_id, user_telegram_id
        FROM customers
        WHERE deleted_at IS NULL AND user_telegram_id = ?
    `

	results, err := db.Query(query, userTelegramID)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	var customer helper.Customer

	for results.Next() {
		err = results.Scan(
			&customer.UserID,
			&customer.UserTelegramID,
		)
		if err != nil {
			panic(err.Error())
		}
	}
	// fmt.Sprintf("UserTelegramId: %d, UserId: %d", &customer.UserTelegramID, &customer.UserID)
	if customer.UserID != nil {
		return &customer
	}
	return nil // тут надо сообщение что если не найден то предложить зарегаться
}
