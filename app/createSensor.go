package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) createSensor(ctx *gin.Context) {
	deviceId := ctx.Param("deviceId")

	if _, err := app.database.GetDevice(deviceId); err != nil {
		ctx.HTML(http.StatusNotFound, "error", gin.H{"message": "Device not found"})
		return
	}

	ctx.HTML(http.StatusOK, "createSensor", gin.H{})
}
