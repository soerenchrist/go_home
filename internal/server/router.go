package server

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
	frontend "github.com/soerenchrist/go_home/internal/app"
	"github.com/soerenchrist/go_home/internal/command"
	"github.com/soerenchrist/go_home/internal/db"
	"github.com/soerenchrist/go_home/internal/device"
	"github.com/soerenchrist/go_home/internal/rules"
	"github.com/soerenchrist/go_home/internal/sensor"
	"github.com/soerenchrist/go_home/internal/value"
)

func NewRouter(database db.Database, outputBindings chan value.SensorValue) *gin.Engine {
	router := gin.Default()

	app := frontend.NewApp(router, database)
	app.ServeHtml()

	devicesController := device.NewDevicesController(database)
	sensorsController := sensor.NewSensorsController(database)
	sensorValuesController := value.NewSensorValuesController(database, outputBindings)
	commandsController := command.NewCommandsController(database)
	rulesController := rules.NewRulesController(database)

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", health)

	v1.GET("/devices", devicesController.GetDevices)
	v1.GET("/devices/:deviceId", devicesController.GetDevice)
	v1.POST("/devices", devicesController.PostDevice)
	v1.DELETE("/devices/:deviceId", devicesController.DeleteDevice)

	v1.GET("/devices/:deviceId/sensors", sensorsController.GetSensors)
	v1.POST("/devices/:deviceId/sensors", sensorsController.PostSensor)
	v1.GET("/devices/:deviceId/sensors/:sensorId", sensorsController.GetSensor)
	v1.DELETE("/devices/:deviceId/sensors/:sensorId", sensorsController.DeleteSensor)

	v1.POST("/devices/:deviceId/sensors/:sensorId/values", sensorValuesController.PostSensorValue)
	v1.GET("/devices/:deviceId/sensors/:sensorId/values", sensorValuesController.GetSensorValues)
	v1.GET("/devices/:deviceId/sensors/:sensorId/current", sensorValuesController.GetCurrentSensorValue)

	v1.GET("/devices/:deviceId/commands", commandsController.GetCommands)
	v1.GET("/devices/:deviceId/commands/:commandId", commandsController.GetCommand)
	v1.POST("/devices/:deviceId/commands", commandsController.PostCommand)
	v1.POST("/devices/:deviceId/commands/:commandId/invoke", commandsController.InvokeCommand)
	v1.DELETE("/devices/:deviceId/commands/:commandId", commandsController.DeleteCommand)

	v1.GET("/rules", rulesController.ListRules)
	v1.POST("/rules", rulesController.PostRule)

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

func health(context *gin.Context) {
	context.JSON(200, gin.H{
		"status": "ok",
	})
}
