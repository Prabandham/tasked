package models

import (
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

// Project model
type Project struct {
	gorm.Model

	Name  string
	Users []User `gorm:many2many:user_projects;`
}

// Validate validates that a project should not be empty
func (project *Project) Validate(v *revel.Validation) {
	v.Required(project.Name)
	v.MinSize(project.Name, 3)
}
