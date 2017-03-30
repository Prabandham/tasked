package models

import (
	"github.com/jinzhu/gorm"
)

type Phase struct {
	gorm.Model

	Name        string
	Description string `gorm:"size:65535"`
	Project     Project
	ProjectID   uint
	Tasks       []Task
}
