package models

import (
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	gorm.Model

	Name                 string
	Email                string
	Password             string
	PasswordConfirmation string `gorm:"-"`
	TermsOfUse           bool

	Projects []Project `gorm:many2many:user_projects;`
	Tasks    []Task
}

// Validate validates the object while registering
func (user *User) Validate(v *revel.Validation) {
	v.Required(user.Name)
	v.MinSize(user.Name, 3)
	v.Required(user.Password)
	v.MinSize(user.Password, 6)
	v.Required(user.PasswordConfirmation)
	v.Required(user.PasswordConfirmation == user.Password).
		Message("The passwords do not match.")
	v.Required(user.Email)
	v.Email(user.Email)
}

// BeforeSave callback to encrypt user password
func (user *User) BeforeSave() (err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashPassword)
	return
}
