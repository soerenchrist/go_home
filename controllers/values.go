package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/mini_home/db"
	"github.com/soerenchrist/mini_home/models"
)

type SensorValuesController struct {
	database db.DevicesDatabase
}

func NewSensorValuesController(database db.DevicesDatabase) *SensorValuesController {
	return &SensorValuesController{database: database}
}

func (c *SensorValuesController) PostSensorValue(context *gin.Context) {
	deviceId := context.Param("deviceId")
	sensorId := context.Param("sensorId")

	var request models.AddSensorValueRequest

	if err := context.BindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sensor, err := c.database.GetSensor(deviceId, sensorId)
	if err != nil {
		context.JSON(404, gin.H{"error": "Sensor not found"})
		return
	}

	if err = c.validateSensorData(&sensor, &request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if request.Timestamp == "" {
		request.Timestamp = getTimestamp()
	}

	sensorValue := models.SensorValue{
		Value:     request.Value,
		Timestamp: request.Timestamp,
		DeviceID:  deviceId,
		SensorID:  sensorId,
	}

	err = c.database.AddSensorValue(sensorValue)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, sensorValue)
}

func getTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

func (c *SensorValuesController) validateSensorData(sensor *models.Sensor, request *models.AddSensorValueRequest) error {
	if sensor.DataType == models.DataTypeInt {
		if _, err := strconv.Atoi(request.Value); err != nil {
			return &models.ValidationError{Message: "Sensor data value is not an int"}
		}
	} else if sensor.DataType == models.DataTypeFloat {
		if _, err := strconv.ParseFloat(request.Value, 64); err != nil {
			return &models.ValidationError{Message: "Sensor data value is not a float"}
		}
	} else if sensor.DataType == models.DataTypeBool {
		if _, err := strconv.ParseBool(request.Value); err != nil {
			return &models.ValidationError{Message: "Sensor data value is not a bool"}
		}
	}

	if sensor.Type == models.SensorTypePolling {
		return &models.ValidationError{Message: "Sensor data value is not allowed for polling sensors"}
	}

	return nil
}
