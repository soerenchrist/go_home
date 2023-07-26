package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) device(ctx *gin.Context) {
	deviceId := ctx.Param("deviceId")
	device, err := app.database.GetDevice(deviceId)
	if err != nil {
		ctx.HTML(http.StatusOK, "not_found", gin.H{"message": "Device not found", "back_link": "/"})
	}
	sensors, err := app.database.ListSensors(deviceId)
	if err != nil {
		ctx.HTML(http.StatusOK, "error", gin.H{})
	}
	commands, err := app.database.ListCommands(deviceId)
	if err != nil {
		ctx.HTML(http.StatusOK, "error", gin.H{})
	}

	ctx.HTML(http.StatusOK, "device", gin.H{
		"device":   device,
		"sensors":  sensors,
		"commands": commands,
	})
}

func (app *App) createDevice(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "createDevice", gin.H{})
}

func (app *App) index(ctx *gin.Context) {
	devices, err := app.database.ListDevices()
	if err != nil {
		panic(err)
	}

	ctx.HTML(http.StatusOK, "index", gin.H{
		"devices": devices,
	})
}
