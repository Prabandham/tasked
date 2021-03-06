package controllers

import (
	"fmt"

	"github.com/revel/revel"

	"github.com/Prabandham/tasked/app"
	"github.com/Prabandham/tasked/app/models"
)

// Dashboard controller
type Dashboard struct {
	Application
}

// Index is going to handle showing of all projects the user has
func (c Dashboard) Index() revel.Result {
	//var user models.User
	//app.DB.First(&user, c.Session["user_id"])

	var projects []models.Project

	query := fmt.Sprintf("SELECT * FROM projects INNER JOIN user_projects ON user_projects.project_id = projects.id WHERE user_projects.user_id = %s", c.Session["user_id"])
	app.DB.Raw(query).Scan(&projects)
	projectsCount := len(projects)

	return c.Render(projects, projectsCount)
}
