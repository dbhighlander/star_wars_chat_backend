// models/chat_message.go
package models

import "time"

type ChatMessage struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	ChatID      uint   `gorm:"not null"`
	BotID       *uint  `gorm:"not null"`
	MessageType string `gorm:"not null"`
	Message     string `gorm:"type:text;not null"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time

	Chat Chat `gorm:"foreignKey:ChatID"`
	Bot  Bot  `gorm:"foreignKey:BotID"`
}
