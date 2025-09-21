package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"first_name" gorm:"size:100;not null"`
	LastName  string `json:"last_name" gorm:"size:100;not null"`
	Email     string `gorm:"size:100;unique;not null"`
	Password  string `json:"password" gorm:"size:255;not null"`
	Role	  string `gorm:"type:enum('user','admin');default:'user'" json:"role"`

	Boards []Board `gorm:"foreignKey:UserID" json:"-"`
	}