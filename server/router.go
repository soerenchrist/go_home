package server

import (
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

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", health.Status)

	v1.GET("/devices", devicesController.GetDevices)
	v1.GET("/devices/:deviceId/", devicesController.GetDevice)
	v1.POST("/devices", devicesController.PostDevice)
	v1.DELETE("/devices/:deviceId/", devicesController.DeleteDevice)

	v1.GET("/devices/:deviceId/sensors", sensorsController.GetSensors)
	v1.POST("/devices/:deviceId/sensors", sensorsController.PostSensor)
	v1.GET("/devices/:deviceId/sensors/:sensorId", sensorsController.GetSensor)
	v1.DELETE("/devices/:deviceId/sensors/:sensorId", sensorsController.DeleteSensor)

	v1.POST("/devices/:deviceId/sensors/:sensorId/values", sensorDataController.PostSensorValue)
	v1.GET("/devices/:deviceId/sensors/:sensorId/values", sensorDataController.GetSensorValues)
	v1.GET("/devices/:deviceId/sensors/:sensorId/current", sensorDataController.GetCurrentSensorValue)

	return router
}
