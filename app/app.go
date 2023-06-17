package app

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/go_home/db"
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
	app.router.HTMLRender = ginview.Default()
	app.router.GET("/", app.index)
	app.router.GET("/devices/:deviceId", app.device)
	app.router.GET("/devices/:deviceId/sensors/:sensorId", app.sensor)
	app.router.GET("/devices/:deviceId/commands/:commandId", app.command)

	app.router.Static("/static", "./views/static")
}
