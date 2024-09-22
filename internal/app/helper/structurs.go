package helper

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
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
