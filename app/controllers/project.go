package controllers

import (
	"strconv"
	"time"

	"github.com/Prabandham/tasked/app"
	"github.com/Prabandham/tasked/app/models"
	"github.com/revel/revel"
)

//Project controller
type Project struct {
	Application
}

// CreateProject handles creation of a project
func (c Project) CreateProject() revel.Result {
	project := models.Project{Name: c.Request.Form.Get("project.Name")}
	project.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/dashboard")
	}
	app.DB.Create(&project)

	intUserID, _ := strconv.Atoi(c.Session["user_id"])

	userProject := models.UserProject{UserId: uint(intUserID), ProjectId: project.ID}
	app.DB.Create(&userProject)

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

// Show handles show page for a project
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

// CreateTask endpoint for creating a task
func (c Project) CreateTask() revel.Result {
	var task models.Task
	phaseID, _ := strconv.Atoi(c.Request.Form.Get("task.PhaseID"))
	userID, _ := strconv.Atoi(c.Session["user_id"])
	task.Name = c.Request.Form.Get("task.Name")
	task.PhaseID = uint(phaseID)
	task.UserID = uint(userID)
	task.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(c.Request.Referer())
	}
	app.DB.Create(&task)

	c.Flash.Success("Created Successfully !!")
	c.Validation.Keep()
	c.FlashParams()

	return c.Redirect(c.Request.Referer())
}

// CompleteTask set's a task as complete for a single user context
func (c Project) CompleteTask() revel.Result {
	var task models.Task
	app.DB.First(&task, c.Request.Form.Get("value"))
	timeNow := time.Now()
	task.CompletedAt = &timeNow
	app.DB.Save(&task)

	response := map[string]string{"status": "success", "cssClass": "list-group-item-success"}
	return c.RenderJson(response)
}

// AddUser add's a user to the specified Project
func (c Project) AddUser() revel.Result {
	var user models.User
	var userProject models.UserProject

	app.DB.Where("email = ?", c.Request.Form.Get("email")).First(&user)
	if user.ID == 0 {
		response := map[string]string{"status": "failure", "message": "No User found by that Email ID"}
		return c.RenderJson(response)
	}
	projectID, _ := strconv.Atoi(c.Request.Form.Get("project_id"))
	app.DB.FirstOrCreate(&userProject, map[string]interface{}{"user_id": user.ID, "project_id": uint(projectID)})
	if userProject.UserId == user.ID {
		c.Flash.Success("You are now Collboarting with %s", user.Name)
		c.Validation.Keep()
		c.FlashParams()
	}
	response := map[string]string{"status": "success", "message": "User added to Project !!"}
	return c.RenderJson(response)
}

// UpdatePhase will handle updating the phase name and Description
func (c Project) UpdatePhase() revel.Result {
	var phase models.Phase

	app.DB.First(&phase, c.Request.Form.Get("phase_id"))
	phase.Name = c.Request.Form.Get("phase_name")
	phase.Description = c.Request.Form.Get("phase_description")
	app.DB.Save(&phase)

	return c.RenderJson(map[string]string{"status": "success", "message": "updated phase"})
}

// UpdateTask Will handle updating the Task
func (c Project) UpdateTask() revel.Result {
	var task models.Task
	formData := c.Request.Form
	app.DB.First(&task, formData.Get("task_id"))

	if formData.Get("task_name") != "" {
		task.Name = formData.Get("task_name")
	}
	if formData.Get("task_description") != "" {
		task.Description = formData.Get("task_description")
	}
	if formData.Get("phase_id") != "" {
		phaseID, _ := strconv.Atoi(formData.Get("phase_id"))
		task.PhaseID = uint(phaseID)
	}
	if formData.Get("user_id") != "" {
		userID, _ := strconv.Atoi(formData.Get("user_id"))
		task.UserID = uint(userID)
	}

	app.DB.Save(&task)
	return c.RenderJson(map[string]string{"status": "success", "message": "updated task"})
}

// AddPhase will add a phase to a given Project
func (c Project) AddPhase() revel.Result {
	var phase models.Phase
	formData := c.Request.Form
	projectID, _ := strconv.Atoi(formData.Get("project_id"))

	app.DB.FirstOrCreate(&phase, map[string]interface{}{"project_id": uint(projectID), "name": formData.Get("phase_name")})
	return c.RenderJson(map[string]string{"status": "success", "message": "created phase"})
}
