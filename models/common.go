package models

import (
	"gorm.io/gorm"
	"time"
)

// auto incremented ID
type ID struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

// Create Time and Update Time
type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Soft Delete
type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
