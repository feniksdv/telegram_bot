package category_sources

import (
	"bot/internal/app/helper"
	"bot/internal/database"
	"database/sql"
	"fmt"
)

// GetAllCategorySourcesByUserId получает все категории для определенного пользователя.
func GetAllCategorySourcesByUserId(userId *int32) ([]helper.CategorySource, error) {
	db := database.Connect()
	defer db.Close()

	query := `
        SELECT 
            id,user_id,name
        FROM category_sources
        WHERE deleted_at IS NULL AND user_id = ?
    `

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categorySources []helper.CategorySource
	for rows.Next() {
		var categorySource helper.CategorySource
		err = rows.Scan(
			&categorySource.ID,
			&categorySource.UserID,
			&categorySource.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		categorySources = append(categorySources, categorySource)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration over result set: %w", err)
	}

	return categorySources, nil
}

// получаем категорию по userId и названию категории
func GetCategorySourceByUserIdAndName(userId *int32, categoryName string) (*helper.CategorySource, error) {
	db := database.Connect()
	defer db.Close()

	query := `
        SELECT 
            id,user_id,name
        FROM category_sources
        WHERE deleted_at IS NULL AND user_id = ? AND name = ?
    `

	if userId == nil {
		return nil, fmt.Errorf("userId is nil")
	}

	var categorySource helper.CategorySource
	err := db.QueryRow(query, *userId, categoryName).Scan(
		&categorySource.ID,
		&categorySource.UserID,
		&categorySource.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Категория %s не найдена для пользователя %d", categoryName, *userId)
		}
		panic(err.Error())
	}

	return &categorySource, nil
}
