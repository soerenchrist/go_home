package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) sensor(ctx *gin.Context) {
	sensorId := ctx.Param("sensorId")
	deviceId := ctx.Param("deviceId")

	sensor, err := app.database.GetSensor(deviceId, sensorId)
	if err != nil {
		ctx.HTML(http.StatusOK, "not_found", gin.H{"message": "Sensor not found", "back_link": "/"})
	}

	ctx.HTML(http.StatusOK, "sensor", gin.H{
		"sensor": sensor,
	})
}
