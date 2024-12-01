package category_expensses

import (
	"bot/internal/app/helper"
	"bot/internal/database"
	"database/sql"
	"fmt"
)

// получаем весь список категорий по id пользователя
func GetCategoriesExpenssesByUserId(userId *int32) ([]helper.CategoryExpenses, error) {
	db := database.Connect()
	defer db.Close()

	query := `
        SELECT 
            id,user_id,name
        FROM category_expenses
        WHERE deleted_at IS NULL AND user_id = ?
    `

	if userId == nil {
		return nil, fmt.Errorf("userId is nil")
	}

	results, err := db.Query(query, *userId)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	var categoryExpensses []helper.CategoryExpenses

	for results.Next() {
		var categoryExpensse helper.CategoryExpenses
		err = results.Scan(
			&categoryExpensse.ID,
			&categoryExpensse.UserID,
			&categoryExpensse.Name,
		)
		if err != nil {
			panic(err.Error())
		}
		categoryExpensses = append(categoryExpensses, categoryExpensse)
	}

	return categoryExpensses, nil // тут надо сообщение что если не найден то предложить зарегаться
}

// получаем категорию по userId и названию категории
func GetCategoryExpensseByUserIdAndName(userId *int32, categoryName string) (*helper.CategoryExpenses, error) {
	db := database.Connect()
	defer db.Close()

	query := `
        SELECT 
            id,user_id,name
        FROM category_expenses
        WHERE deleted_at IS NULL AND user_id = ? AND name = ?
    `

	if userId == nil {
		return nil, fmt.Errorf("userId is nil")
	}

	var categoryExpensse helper.CategoryExpenses
	err := db.QueryRow(query, *userId, categoryName).Scan(
		&categoryExpensse.ID,
		&categoryExpensse.UserID,
		&categoryExpensse.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Категория %s не найдена для пользователя %d", categoryName, *userId)
		}
		panic(err.Error())
	}

	return &categoryExpensse, nil
}
