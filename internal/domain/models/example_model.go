package models

import (
	"time"
)

type Example struct {
	Id        int        `json:"id" gorm:"primaryKey"`
	Data      string     `json:"data" binding:"required"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at_at"`
}

// TableName is Database TableName of this model
func (e *Example) TableName() string {
	return "examples"
}
