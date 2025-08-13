package models

import (
	"time"
)

type Task struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	Completed bool      `json:"completed" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	Deadline  time.Time `json:"deadline"`

	UserID int  `json:"user_id" gorm:"column:user_id;not null"`
	User   User `json:"-" gorm:"foreignKey:UserID;references:ID"`
}

type User struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserName    string `json:"user_name" gorm:"type:varchar(255);not null"`
	Password    string `json:"-" gorm:"type:varchar(255);not null"`
	PhoneNumber string `json:"phone_number" gorm:"not null"`
}
