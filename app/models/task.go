package models

import (
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"time"
)

type Task struct {
	gorm.Model

	Name        string
	Description string `gorm:"size:65535"`
	Phase       Phase
	PhaseID     uint
	User        User
	UserID      uint
	CompletedAt *time.Time
	DueOn       *time.Time
}

func (task *Task) Validate(v *revel.Validation) {
	v.Required(task.Name)
	v.MinSize(task.Name, 3)
}

func (t Task) IsComplete() bool {
	if t.CompletedAt == nil {
		return false
	}
	return true
}

func (t Task) CompletionStatus() string {
	if t.IsComplete() == true {
		return "list-group-item-success"
	}
	return "list-group-item-warning"
}
