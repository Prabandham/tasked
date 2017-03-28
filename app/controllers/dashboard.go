package controllers

import (
	"github.com/revel/revel"

	"fmt"
	"github.com/Prabandham/tasked/app"
	"github.com/Prabandham/tasked/app/models"
	"strconv"
)

// Dashboard controller
type Dashboard struct {
	Application
}

func (c Dashboard) Index() revel.Result {
	var user models.User
	app.DB.First(&user, c.Session["user_id"])

	var projects []models.Project

	query := fmt.Sprintf("SELECT * FROM projects INNER JOIN user_projects ON user_projects.project_id = projects.id WHERE user_projects.user_id = %d", user.ID)
	app.DB.Raw(query).Scan(&projects)
	projectsCount := len(projects)

	return c.Render(projects, projectsCount)
}

func (c Dashboard) CreateProject() revel.Result {
	project := models.Project{Name: c.Request.Form.Get("project.Name")}
	app.DB.Create(&project)

	intUserId, _ := strconv.Atoi(c.Session["user_id"])

	userProject := models.UserProject{UserId: uint(intUserId), ProjectId: project.ID}
	app.DB.Create(&userProject)

	c.Flash.Success("Created Successfully !!")
	c.Validation.Keep()
	c.FlashParams()

	return c.Redirect("/dashboard")
}
