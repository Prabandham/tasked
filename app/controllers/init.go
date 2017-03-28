package controllers

import (
	"github.com/revel/revel"
)

func init() {
	revel.InterceptFunc(CheckLoggedInUser, revel.BEFORE, &Application{})
}

// CheckLoggedInUser checks to see if the user has a valid session and then allows him
// to access the Application
func CheckLoggedInUser(c *revel.Controller) revel.Result {

	// IF the instance is that of an Application then we check for logged in user
	if c.Action == "App.Login" || c.Action == "App.Register" || c.Action == "App.Index" || c.Action == "Static.Serve" {
		// Do nothing just Continue
	} else {
		if c.Session["user_name"] == "" && c.Session["user_id"] == "" {
			c.Flash.Error("Please Login to Continue !!")
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect("/login")
		}
		//TODO query the user and then proceed
	}
	return nil
}
