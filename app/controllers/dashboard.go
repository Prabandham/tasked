package controllers

import (
	"github.com/revel/revel"
)

// Dashboard controller
type Dashboard struct {
	Application
}

func (c Dashboard) Index() revel.Result {
	return c.Render()
}
