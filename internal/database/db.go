package database

import (
	"bot/internal/app/helper"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// TODO перенос в env
func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3309)/microfo")
	if err != nil {
		log.Print(err.Error())
	}

	return db
}

func GetByUserTelegramId() *string {
	db := connect()
	defer db.Close()
	results, err := db.Query("SELECT user_telegram_id, user_id, deleted_at FROM session_bots WHERE deleted_at IS NULL")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var tag helper.SessionBots

	for results.Next() {
		// для каждой строки отсканируйте результат в нашем теге composite object
		err = results.Scan(&tag.UserTelegramId, &tag.UserId, &tag.DeletedAt)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	if tag.UserId != nil {
		return tag.UserId // Возвращаем указатель на UserId
	}
	return nil // Возвращаем nil, если UserId не найден
}

func GetFirst() *string {
	db := connect()
	defer db.Close()
	results, err := db.Query("SELECT user_telegram_id, user_id FROM session_bots")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var tag helper.SessionBots

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.UserTelegramId, &tag.UserId)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	// and then print out the tag's Name attribute
	return tag.UserId
}

func Update() {

}

func Insert() {

}

func Delete() {

}
