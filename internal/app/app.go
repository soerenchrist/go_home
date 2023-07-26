package app

import (
	"html/template"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/go_home/internal/db"
)

type App struct {
	router   *gin.Engine
	database db.Database
}

func NewApp(router *gin.Engine, database db.Database) *App {
	return &App{
		router:   router,
		database: database,
	}
}

func (app *App) ServeHtml() {
	view := ginview.New(goview.Config{
		Root:         "internal/app/views",
		Extension:    ".html",
		Master:       "layouts/master",
		Partials:     []string{},
		Funcs:        make(template.FuncMap),
		DisableCache: false,
		Delims:       goview.Delims{Left: "{{", Right: "}}"},
	})
	app.router.HTMLRender = view
	app.router.GET("/", app.index)
	app.router.GET("/devices/:deviceId", app.device)
	app.router.GET("/devices/:deviceId/sensors/:sensorId", app.sensor)
	app.router.GET("/devices/:deviceId/commands/:commandId", app.command)

	app.router.GET("/createDevice", app.createDevice)
	app.router.GET("/devices/:deviceId/createSensor", app.createSensor)
	app.router.GET("/devices/:deviceId/createCommand", app.createCommand)

	app.router.GET("/createRule", app.createRule)

	app.router.Static("/static", "./internal/app/views/static")
}
