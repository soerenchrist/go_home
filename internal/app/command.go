package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) command(ctx *gin.Context) {
	commandId := ctx.Param("commandId")
	deviceId := ctx.Param("deviceId")

	command, err := app.database.GetCommand(deviceId, commandId)
	if err != nil {
		ctx.HTML(http.StatusOK, "not_found", gin.H{"message": "Command not found", "back_link": "/"})
	}

	ctx.HTML(http.StatusOK, "command", gin.H{
		"command": command,
	})
}

func (app *App) createCommand(ctx *gin.Context) {
	deviceId := ctx.Param("deviceId")

	if _, err := app.database.GetDevice(deviceId); err != nil {
		ctx.HTML(http.StatusNotFound, "error", gin.H{"message": "Device not found"})
		return
	}

	ctx.HTML(http.StatusOK, "createCommand", gin.H{})
}
