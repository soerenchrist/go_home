package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/soerenchrist/mini_home/db"
	"github.com/soerenchrist/mini_home/models"
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

func (c *SensorsController) PostSensor(context *gin.Context) {
	deviceId := context.Param("deviceId")
	var request models.CreateSensorRequest
	if err := context.BindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.validateSensor(request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sensor := models.Sensor{
		ID:       uuid.NewString(),
		Name:     request.Name,
		DeviceID: deviceId,
		DataType: models.DataType(request.DataType),
		Unit:     request.Unit,
	}

	err := c.database.AddSensor(sensor)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, sensor)
}

func (c *SensorsController) validateSensor(sensor models.CreateSensorRequest) error {
	if len(sensor.Name) < 3 {
		return &models.ValidationError{Message: "Name must be at least 3 characters long"}
	}

	if len(sensor.Unit) > 0 && (sensor.DataType == models.DataTypeBool || sensor.DataType == models.DataTypeString) {
		return &models.ValidationError{Message: "Unit is not allowed for this data type"}
	}

	return nil
}
