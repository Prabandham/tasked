package controllers

import (
	"github.com/Prabandham/tasked/app"
	"github.com/Prabandham/tasked/app/models"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

// App Controller
type App struct {
	*revel.Controller
}

// Index is the default route that is accessible at /
func (c App) Index() revel.Result {
	return c.Render()
}

// Login will handle the login for the application
func (c App) Login() revel.Result {
	if c.Request.Method == "POST" {
		var user models.User
		password := []byte(c.Params.Form.Get("password"))
		hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
		if err != nil {
			panic(err)
		}
		app.DB.Find(&user, "email = ? and password = ?", c.Params.Form.Get("email"), hashedPassword)
		if user.Email != "" {
			//The user exists put this guy in session
		}
		c.Flash.Error("Invalid Email or Password !!")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Login)
	}
	return c.Render()
}

// Register will handle registering the users
func (c App) Register(user *models.User) revel.Result {
	if c.Request.Method == "POST" {
		c.Validation.Clear()
		user.Validate(c.Validation)
		if c.Validation.HasErrors() {
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect(App.Register)
		}
		if app.DB.NewRecord(user) {
			app.DB.Create(user)
			if user.ID == 0 {
				c.Flash.Error("Email already exists !! Please login to continue !!")
				c.Validation.Keep()
				c.FlashParams()
				return c.Redirect(App.Login)
			}
		}
	}
	return c.Render()
}
