package devices

import "github.com/gin-gonic/gin"

var devices = []Device{
	{ID: "1", Name: "Device 1", LastReached: "2018-01-01T00:00:00Z"},
}

func getDevices(context *gin.Context) {
	context.JSON(200, devices)
}

func MapEndpoints(router *gin.Engine) {
	router.GET("/api/devices", getDevices)
}
