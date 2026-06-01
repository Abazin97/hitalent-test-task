package models

import "time"

type Department struct {
	ID        int64  `gorm:"primaryKey"`
	Name      string `gorm:"size:200;not null"`
	ParentID  *int64
	CreatedAt time.Time

	Parent   *Department  `gorm:"foreignKey:ParentID"`
	Children []Department `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`

	Employees []Employee `gorm:"foreignKey:DepartmentID"`
}
