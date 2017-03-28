package controllers

import (
	"github.com/revel/revel"
)

// This is going to server as the main base controller for all logged in controllers.
type Application struct {
	*revel.Controller
}
