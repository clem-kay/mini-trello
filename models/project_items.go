package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectItem struct {
	gorm.Model
	Name        string       `gorm:"size:50;not null" json:"name"`
	BoardID     uint         `gorm:"not null;index" json:"board_id"`
	Description string       `gorm:"size:255" json:"description"`
	DueDate     *time.Time   `json:"due_date,omitempty"` // optional
	Status      ItemStatus   `gorm:"type:varchar(20);default:'todo';index" json:"status"` // safer for cross-db
	Priority    ItemPriority `gorm:"type:varchar(10);default:'medium';index" json:"priority"`

	Board *Board `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE" json:"-"`
}

type ItemStatus string

const (
	StatusTodo       ItemStatus = "todo"
	StatusInProgress ItemStatus = "in_progress"
	StatusDone       ItemStatus = "done"
)

type ItemPriority string

const (
	PriorityLow    ItemPriority = "low"
	PriorityMedium ItemPriority = "medium"
	PriorityHigh   ItemPriority = "high"
)