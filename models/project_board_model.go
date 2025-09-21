// models/board.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Board struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	UserID      uint     `gorm:"not null;index" json:"user_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
}
