package models

import "time"

type Bot struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string     `gorm:"type:varchar(120);not null" json:"name"`
	Slug              string     `gorm:"type:varchar(120);not null" json:"slug"`
	PersonalityPrompt string     `gorm:"type:text;not null" json:"personality_prompt"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"modified_at"`
}
