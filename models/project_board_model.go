package models

import "gorm.io/gorm"

type Board struct {
	gorm.Model
	Name        string        `gorm:"size:100;not null" json:"name"`
	Description string        `gorm:"size:255" json:"description"`
	UserID      uint          `gorm:"not null;index" json:"user_id"`

	User  *User         `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
	Items []ProjectItem `gorm:"foreignKey:BoardID" json:"-"` // optional
}