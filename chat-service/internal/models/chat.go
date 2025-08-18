// models/chat.go
package models

import "time"

type Chat struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Reference string `gorm:"type:varchar(255);not null"`
	UserID    uint   `gorm:"not null"`
	BotID     *uint
	CreatedAt *time.Time
	UpdatedAt *time.Time

	User         User          `gorm:"foreignKey:UserID"`
	Bot          *Bot          `gorm:"foreignKey:BotID"`
	ChatMessages []ChatMessage `gorm:"foreignKey:ChatID"`
}
