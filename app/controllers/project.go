package controllers

import (
	"github.com/Prabandham/tasked/app"
	"github.com/Prabandham/tasked/app/models"
	"github.com/revel/revel"
	"strconv"
	"time"
)

type Project struct {
	Application
}

func (c Project) CreateProject() revel.Result {
	project := models.Project{Name: c.Request.Form.Get("project.Name")}
	app.DB.Create(&project)

	// This can be running as a goroutine as it wont affect any thing
	// Once we have the project this should be just a chore
	// TODO may be move this to a afterCreate hook on project
	go func() {
		intUserId, _ := strconv.Atoi(c.Session["user_id"])

		userProject := models.UserProject{UserId: uint(intUserId), ProjectId: project.ID}
		app.DB.Create(&userProject)
	}()

	// This can be running as a goroutine as it wont affect any thing
	// Once we have the project this should be just a chore
	// TODO may be move this to a afterCreate hook on project
	go func() {
		defaultPhase := models.Phase{
			Name:      "Task List",
			ProjectID: project.ID,
		}
		app.DB.Create(&defaultPhase)
	}()

	c.Flash.Success("Created Successfully !!")
	c.Validation.Keep()
	c.FlashParams()

	return c.Redirect("/dashboard")
}

func (c Project) Show(id uint) revel.Result {
	var project models.Project
	var phases []models.Phase
	var tasks []models.Task
	var phaseIds []interface{}
	var projectUsers []interface{}

	app.DB.First(&project, id).Related(&phases)
	// Get a list of PhaseIDs from phases
	for _, phase := range phases {
		phaseIds = append(phaseIds, phase.ID)
	}
	app.DB.Model(models.Task{}).Where("phase_id in (?)", phaseIds).Scan(&tasks)
	app.DB.Raw("SELECT * FROM user_projects WHERE project_id = ?", project.ID).Scan(&projectUsers)
	projectUsersCount := int(len(projectUsers))

	return c.Render(project, phases, tasks, projectUsersCount)
}

func (c Project) CreateTask() revel.Result {
	var task models.Task
	phaseId, _ := strconv.Atoi(c.Request.Form.Get("task.PhaseID"))
	userId, _ := strconv.Atoi(c.Session["user_id"])
	task.Name = c.Request.Form.Get("task.Name")
	task.PhaseID = uint(phaseId)
	task.UserID = uint(userId)
	app.DB.Create(&task)

	c.Flash.Success("Created Successfully !!")
	c.Validation.Keep()
	c.FlashParams()

	return c.Redirect(c.Request.Referer())
}

func (c Project) CompleteTask() revel.Result {
	var task models.Task
	app.DB.First(&task, c.Request.Form.Get("value"))
	timeNow := time.Now()
	task.CompletedAt = &timeNow
	app.DB.Save(&task)

	response := map[string]string{"status": "success", "cssClass": "list-group-item-success"}
	return c.RenderJson(response)
}
