package models

import "time"

type Employee struct {
	ID           int64 `gorm:"primaryKey"`
	DepartmentID int64 `gorm:"not null"`

	FullName string `gorm:"size:200;not null"`
	Position string `gorm:"size:200;not null"`

	HiredAt   *time.Time
	CreatedAt time.Time

	Department Department
}
