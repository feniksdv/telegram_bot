package category_incomes

import (
	"bot/internal/app/helper"
	"bot/internal/database"
	"database/sql"
	"fmt"
)

// получаем весь список категорий по id пользователя
func GetCategoriesIncomesByUserId(userId *int32) ([]helper.CategoryIncomes, error) {
	db := database.Connect()
	defer db.Close()

	query := `
        SELECT 
            id,user_id,name
        FROM category_incomes
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

	var categoryIncomes []helper.CategoryIncomes

	for results.Next() {
		var categoryIncome helper.CategoryIncomes
		err = results.Scan(
			&categoryIncome.Id,
			&categoryIncome.UserId,
			&categoryIncome.Name,
		)
		if err != nil {
			panic(err.Error())
		}
		categoryIncomes = append(categoryIncomes, categoryIncome)
	}

	return categoryIncomes, nil // тут надо сообщение что если не найден то предложить зарегаться
}

// получаем категорию по userId и названию категории
func GetCategoryIncomeByUserIdAndName(userId *int32, categoryName string) (*helper.CategoryIncomes, error) {
	db := database.Connect()
	defer db.Close()

	query := `
        SELECT 
            id,user_id,name
        FROM category_incomes
        WHERE deleted_at IS NULL AND user_id = ? AND name = ?
    `

	if userId == nil {
		return nil, fmt.Errorf("userId is nil")
	}

	var categoryIncome helper.CategoryIncomes
	err := db.QueryRow(query, *userId, categoryName).Scan(
		&categoryIncome.Id,
		&categoryIncome.UserId,
		&categoryIncome.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Категория %s не найдена для пользователя %d", categoryName, *userId)
		}
		panic(err.Error())
	}

	return &categoryIncome, nil
}
