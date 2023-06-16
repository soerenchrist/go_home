package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
	"github.com/soerenchrist/go_home/util"
)

type SensorValuesController struct {
	database       db.DevicesDatabase
	outputBindings chan models.SensorValue
}

func NewSensorValuesController(database db.DevicesDatabase, outputBindings chan models.SensorValue) *SensorValuesController {
	return &SensorValuesController{database: database, outputBindings: outputBindings}
}

func (c *SensorValuesController) GetSensorValues(context *gin.Context) {
	sensor, device, err := c.getSensorAndDevice(context)
	if err != nil {
		context.JSON(404, gin.H{"error": err.Error()})
		return
	}

	timeframeQuery, isOk := context.GetQuery("timeframe")
	if !isOk {
		timeframeQuery = "1h"
	}

	timeframe, err := time.ParseDuration(timeframeQuery)
	if err != nil {
		context.JSON(400, gin.H{"error": "Invalid timeframe"})
		return
	}

	since := time.Now().Add(-timeframe)

	sensorValues, err := c.database.GetSensorValuesSince(device.ID, sensor.ID, since)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, sensorValues)
}

func (c *SensorValuesController) GetCurrentSensorValue(context *gin.Context) {
	sensor, device, err := c.getSensorAndDevice(context)
	if err != nil {
		context.JSON(404, gin.H{"error": err.Error()})
		return
	}

	sensorValue, err := c.database.GetCurrentSensorValue(device.ID, sensor.ID)
	if err != nil {
		log.Println(err)
		context.JSON(404, gin.H{"error": "No sensor value found"})
		return
	}

	context.JSON(200, sensorValue)
}

func (c *SensorValuesController) PostSensorValue(context *gin.Context) {
	sensor, device, err := c.getSensorAndDevice(context)
	if err != nil {
		context.JSON(404, gin.H{"error": err.Error()})
		return
	}

	var request models.AddSensorValueRequest

	if err := context.BindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err = c.validateSensorData(sensor, &request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if request.Timestamp == "" {
		request.Timestamp = util.GetTimestamp()
	}

	sensorValue := &models.SensorValue{
		Value:     request.Value,
		Timestamp: request.Timestamp,
		DeviceID:  device.ID,
		SensorID:  sensor.ID,
	}

	err = c.database.AddSensorValue(sensorValue)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.outputBindings <- *sensorValue
	context.JSON(201, sensorValue)
}

func (c *SensorValuesController) getSensorAndDevice(context *gin.Context) (*models.Sensor, *models.Device, error) {
	deviceId := context.Param("deviceId")
	sensorId := context.Param("sensorId")

	var device *models.Device
	var sensor *models.Sensor

	var err error

	if device, err = c.database.GetDevice(deviceId); err != nil {
		return &models.Sensor{}, &models.Device{}, &models.NotFoundError{Message: "Device not found"}
	}

	if sensor, err = c.database.GetSensor(deviceId, sensorId); err != nil {
		return &models.Sensor{}, &models.Device{}, &models.NotFoundError{Message: "Sensor not found"}
	}

	return sensor, device, nil
}

func (c *SensorValuesController) validateSensorData(sensor *models.Sensor, request *models.AddSensorValueRequest) error {
	if sensor.DataType == models.DataTypeInt {
		if _, err := strconv.Atoi(request.Value); err != nil {
			return &models.ValidationError{Message: "Sensor value is not an int"}
		}
	} else if sensor.DataType == models.DataTypeFloat {
		if _, err := strconv.ParseFloat(request.Value, 64); err != nil {
			return &models.ValidationError{Message: "Sensor value is not a float"}
		}
	} else if sensor.DataType == models.DataTypeBool {
		if _, err := strconv.ParseBool(request.Value); err != nil {
			return &models.ValidationError{Message: "Sensor value is not a bool"}
		}
	}

	if sensor.Type == models.SensorTypePolling {
		return &models.ValidationError{Message: "Sending values to a polling sensor is not allowed"}
	}

	if !sensor.IsActive {
		return &models.ValidationError{Message: "Sensor is not active"}
	}

	if len(request.Timestamp) > 0 {
		if err := util.ValidateTimestamp(request.Timestamp); err != nil {
			return &models.ValidationError{Message: fmt.Sprintf("%s is not a valid RFC3339 timestamp", request.Timestamp)}
		}
	}

	return nil
}
