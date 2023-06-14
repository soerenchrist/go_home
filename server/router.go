package server

import (
	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/mini_home/controllers"
	"github.com/soerenchrist/mini_home/db"
)

func NewRouter(database db.DevicesDatabase) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)
	devicesController := controllers.NewDevicesController(database)
	sensorsController := controllers.NewSensorsController(database)

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", health.Status)

	v1.GET("/devices", devicesController.GetDevices)
	v1.GET("/devices/:id/", devicesController.GetDevice)
	v1.POST("/devices", devicesController.PostDevice)
	v1.DELETE("/devices/:id/", devicesController.DeleteDevice)

	v1.GET("/devices/:id/sensors", sensorsController.GetSensors)
	v1.POST("/devices/:id/sensors", sensorsController.PostSensor)

	return router
}
