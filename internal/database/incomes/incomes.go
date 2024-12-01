package incomes

import "bot/internal/database"

func Create() {
	db := database.Connect()
	defer db.Close()
	results, err := db.Prepare("INSERT INTO session_bots(user_telegram_id, user_id, money, money_message_id, typy, unit_id, step) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	results.Exec('1', '1', '1', '1', '1', '1', '1', '1')
}
