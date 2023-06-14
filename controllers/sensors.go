package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/mini_home/db"
)

type SensorsController struct {
	database db.DevicesDatabase
}

func NewSensorsController(database db.DevicesDatabase) *SensorsController {
	return &SensorsController{database: database}
}

func (c *SensorsController) GetSensors(context *gin.Context) {
	deviceId := context.Param("deviceId")

	sensors, err := c.database.ListSensors(deviceId)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, sensors)
}
