package helper

import (
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Contact struct {
	FirstName string
	LastName  string
	UserName  string
}

func GetContactName(update tgbotapi.Update) string {
	contact := Contact{
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		UserName:  update.Message.From.UserName,
	}
	var contactName string
	names := []string{contact.FirstName, contact.LastName, contact.UserName}
	for _, name := range names {
		if name := strings.TrimSpace(name); name != "" {
			contactName = name
			break
		}
	}

	return contactName
}

type SessionBots struct {
	UserTelegramId       *int64     `json:"user_telegram_id"`
	UserId               *int32     `json:"user_id"`
	Money                *string    `json:"money"`
	MoneyMessageId       *int       `json:"money_message_id"`
	CategoryId           *uint64    `json:"category_id"`
	CategoryName         *string    `json:"category_name"`
	CategoryMessageId    *uint64    `json:"category_message_id"`
	Type                 *string    `json:"type"`
	UnitId               *uint8     `json:"unit_id"`
	Description          *string    `json:"description"`
	DescriptionMessageId *uint64    `json:"description_message_id"`
	CategorySourceId     *uint64    `json:"category_source_id"`
	CategorySourceName   *string    `json:"category_source_name"`
	SourceMessageId      *uint64    `json:"source_message_id"`
	Step                 *uint8     `json:"step"`
	IsEdit               *bool      `json:"is_edit"`
	CreatedAt            *time.Time `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at"`
}

type Customer struct {
	ID             *uint64    `json:"id"`
	UserID         *int32     `json:"user_id"`
	UserTelegramID *int64     `json:"user_telegram_id,omitempty"`
	Username       string     `json:"username,omitempty"`
	SendEmail      bool       `json:"send_email"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type CategoryIncomes struct {
	Id          *uint64    `json:"id"`
	UserId      *uint32    `json:"user_id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description,omitempty"` // Новое поле для описания
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type CategorySource struct {
	ID          *uint64    `json:"id"`
	UserID      *uint64    `json:"user_id"`
	Name        *string    `json:"name" validate:"required"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type CategoryExpenses struct {
	ID          *uint64    `json:"id"`
	UserID      *uint64    `json:"user_id"`
	Name        *string    `json:"name" validate:"required"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
