package server

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/mini_home/controllers"
	"github.com/soerenchrist/mini_home/db"
	"github.com/soerenchrist/mini_home/web"
)

func NewRouter(database db.DevicesDatabase) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	web.ServeHtml(router)

	health := new(controllers.HealthController)
	devicesController := controllers.NewDevicesController(database)
	sensorsController := controllers.NewSensorsController(database)
	sensorDataController := controllers.NewSensorValuesController(database)
	commandsController := controllers.NewCommandsController(database)

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", health.Status)

	v1.GET("/devices", devicesController.GetDevices)
	v1.GET("/devices/:deviceId", devicesController.GetDevice)
	v1.POST("/devices", devicesController.PostDevice)
	v1.DELETE("/devices/:deviceId", devicesController.DeleteDevice)

	v1.GET("/devices/:deviceId/sensors", sensorsController.GetSensors)
	v1.POST("/devices/:deviceId/sensors", sensorsController.PostSensor)
	v1.GET("/devices/:deviceId/sensors/:sensorId", sensorsController.GetSensor)
	v1.DELETE("/devices/:deviceId/sensors/:sensorId", sensorsController.DeleteSensor)

	v1.POST("/devices/:deviceId/sensors/:sensorId/values", sensorDataController.PostSensorValue)
	v1.GET("/devices/:deviceId/sensors/:sensorId/values", sensorDataController.GetSensorValues)
	v1.GET("/devices/:deviceId/sensors/:sensorId/current", sensorDataController.GetCurrentSensorValue)

	v1.GET("/devices/:deviceId/commands", commandsController.GetCommands)
	v1.GET("/devices/:deviceId/commands/:commandId", commandsController.GetCommand)
	v1.POST("/devices/:deviceId/commands", commandsController.PostCommand)
	v1.POST("/devices/:deviceId/commands/:commandId/invoke", commandsController.InvokeCommand)
	v1.DELETE("/devices/:deviceId/commands/:commandId", commandsController.DeleteCommand)

	router.POST("/echo", echo)
	return router
}

func echo(context *gin.Context) {
	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.Status(200)
	context.Header("Content-Type", "application/json")
	context.Writer.Write(body)
}
