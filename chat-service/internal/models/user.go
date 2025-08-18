// models/user.go
package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Reference string `gorm:"type:varchar(120);not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time

	Chats []Chat `gorm:"foreignKey:UserID"`
}
