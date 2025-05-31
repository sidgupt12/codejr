package models

import "time"


type Note struct {
	ID int `gorm:"primaryKey" json:"id"`
	Title string `gorm:"not null" json:"title"`
	Content string `gorm:"not null" json:"content"`
	UserId int `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}