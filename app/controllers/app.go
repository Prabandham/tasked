package controllers

import (
	"github.com/Prabandham/tasked/app"
	"github.com/Prabandham/tasked/app/models"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"strconv"
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

		app.DB.Find(&user, "email = ?", c.Params.Form.Get("email"))
		password := []byte(c.Params.Form.Get("password"))
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), password)

		if (user.Name == "") || (err != nil) {
			c.Flash.Error("Invalid Email or Password !!")
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect(App.Login)
		}
		// This has to be enctypted and saved
		c.Session["user_name"] = user.Name
		c.Session["user_id"] = strconv.Itoa(int(user.ID))
		c.Session.SetDefaultExpiration()

		c.Flash.Success("Welcome back, %s", user.Name)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Dashboard.Index)
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
			} else {
				c.Flash.Success("Registerd successfully !! Welcome %s", user.Name)
				c.Validation.Keep()
				c.FlashParams()
				// TODO sign the user in and redirect to dashboard
				return c.Redirect(App.Login)
			}
		}
	}
	return c.Render()
}

// Logout will clear the user Session
func (c App) Logout() revel.Result {
	delete(c.Session, "user_name")
	delete(c.Session, "user_id")
	c.Flash.Success("Logged Out Successfully !!")
	c.Validation.Keep()
	c.FlashParams()
	return c.Redirect(App.Login)
}
