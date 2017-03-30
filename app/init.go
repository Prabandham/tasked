package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // Mysql driver

	"github.com/revel/revel"

	"github.com/Prabandham/tasked/app/models"
)

// DB is the global database object, accebile via app.DB
var DB *gorm.DB

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	revel.OnAppStart(InitDB)
	revel.OnAppStart(MigrateDB)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

// InitDB will create a connection to the database
func InitDB() {
	var err error
	DB, err = gorm.Open("mysql", "root:root@/tasked?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic(err)
	}

	DB.LogMode(true)
}

// MigrateDB Will  add the necessary migrations.
func MigrateDB() {

	if DB.HasTable(&models.User{}) != true {
		DB.CreateTable(&models.User{})
		DB.Model(&models.User{}).AddIndex("idx_user_name", "name")
		DB.Model(&models.User{}).AddUniqueIndex("idx_user_email", "email")
		DB.Model(&models.User{}).AddIndex("idx_user_delete_at", "deleted_at")
		DB.Model(&models.User{}).AddIndex("idx_user_created_at", "created_at")
	}

	if DB.HasTable(&models.Project{}) != true {
		DB.CreateTable(&models.Project{})
		DB.Model(&models.Project{}).AddIndex("idx_project_name", "name")
		DB.Model(&models.Project{}).AddIndex("idx_project_delete_at", "deleted_at")
		DB.Model(&models.Project{}).AddIndex("idx_project_created_at", "created_at")
		DB.Model(&models.Project{}).AddIndex("idx_project_updated_at", "updated_at")
	}

	if DB.HasTable(&models.UserProject{}) != true {
		DB.CreateTable(&models.UserProject{})
		DB.Model(&models.UserProject{}).AddIndex("idx_user_join_id", "user_id")
		DB.Model(&models.UserProject{}).AddIndex("idx_project_join_id", "project_id")
	}

	if DB.HasTable(&models.Phase{}) != true {
		DB.CreateTable(&models.Phase{})
		DB.Model(&models.Phase{}).AddIndex("idx_phase_name", "name")
		DB.Model(&models.Phase{}).AddIndex("idx_phase_project_id", "project_id")
		DB.Model(&models.Phase{}).AddIndex("idx_phase_deleted_at", "deleted_at")
		DB.Model(&models.Phase{}).AddIndex("idx_phase_created_at", "created_at")
		DB.Model(&models.Phase{}).AddIndex("idx_phase_updated_at", "updated_at")
	}

	if DB.HasTable(&models.Task{}) != true {
		DB.CreateTable(&models.Task{}).AddIndex("idx_task_name", "name")
		DB.Model(&models.Task{}).AddIndex("idx_task_phase", "phase_id")
		DB.Model(&models.Task{}).AddIndex("idx_task_user", "user_id")
		DB.Model(&models.Task{}).AddIndex("idx_task_due_on", "due_on")
		DB.Model(&models.Task{}).AddIndex("idx_task_deleted_at", "deleted_at")
		DB.Model(&models.Task{}).AddIndex("idx_task_created_at", "created_at")
		DB.Model(&models.Task{}).AddIndex("idx_task_updated_at", "updated_at")
	}
}
