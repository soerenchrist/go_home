package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/go_home/internal/models"
)

type SensorsDatabase interface {
	ListSensors(deviceId string) ([]models.Sensor, error)
	GetSensor(deviceId string, sensorId string) (*models.Sensor, error)
	GetDevice(deviceId string) (*models.Device, error)
	AddSensor(sensor *models.Sensor) error
	DeleteSensor(deviceId string, sensorId string) error
}

type SensorsController struct {
	database SensorsDatabase
}

func NewSensorsController(database SensorsDatabase) *SensorsController {
	return &SensorsController{database: database}
}

func (c *SensorsController) GetSensors(context *gin.Context) {
	deviceId := context.Param("deviceId")

	if _, err := c.database.GetDevice(deviceId); err != nil {
		context.JSON(404, gin.H{"error": "Device not found"})
		return
	}

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
		context.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if _, err := c.database.GetDevice(deviceId); err != nil {
		context.JSON(404, gin.H{"error": "Device not found"})
		return
	}

	if _, err := c.database.GetSensor(deviceId, request.Id); err == nil {
		context.JSON(400, gin.H{"error": fmt.Sprintf("Sensor with id %s does already exist", request.Id)})
		return
	}

	if err := c.validateSensor(request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if request.Type != models.SensorTypePolling {
		request.Type = models.SensorTypeExternal
	}

	if request.Type == models.SensorTypePolling {
		request.PollingStrategy = models.PollingStrategyPing
	}

	sensor := &models.Sensor{
		ID:              request.Id,
		Name:            request.Name,
		DeviceID:        deviceId,
		DataType:        models.DataType(request.DataType),
		Unit:            request.Unit,
		Type:            request.Type,
		PollingInterval: request.PollingInterval,
		PollingEndpoint: request.PollingEndpoint,
		PollingStrategy: request.PollingStrategy,
		IsActive:        true,
	}

	err := c.database.AddSensor(sensor)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(201, sensor)
}

func (c *SensorsController) GetSensor(context *gin.Context) {
	deviceId := context.Param("deviceId")
	sensorId := context.Param("sensorId")

	device, err := c.database.GetSensor(deviceId, sensorId)
	if err != nil {
		context.JSON(404, gin.H{"error": "Sensor not found"})
		return
	}
	context.JSON(200, device)
}

func (c *SensorsController) DeleteSensor(context *gin.Context) {
	deviceId := context.Param("deviceId")
	sensorId := context.Param("sensorId")

	err := c.database.DeleteSensor(deviceId, sensorId)

	if notFound, isOk := err.(*models.NotFoundError); isOk {
		context.JSON(404, gin.H{"error": notFound.Error()})
		return
	}

	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.Status(204)
}

func (c *SensorsController) validateSensor(sensor models.CreateSensorRequest) error {
	if len(sensor.Name) < 3 {
		return &models.ValidationError{Message: "Name must be at least 3 characters long"}
	}

	if len(sensor.Unit) > 0 && (sensor.DataType == models.DataTypeBool || sensor.DataType == models.DataTypeString) {
		return &models.ValidationError{Message: "Unit is not allowed for this data type"}
	}

	if sensor.Type != models.SensorTypePolling && sensor.Type != models.SensorTypeExternal {
		return &models.ValidationError{Message: "Invalid sensor type"}
	}

	if sensor.Type == models.SensorTypePolling && sensor.PollingInterval < 1 {
		return &models.ValidationError{Message: "Polling interval must be greater than 0"}
	}

	if sensor.Type == models.SensorTypePolling {
		if sensor.PollingStrategy != models.PollingStrategyPing {
			return &models.ValidationError{Message: "Invalid polling strategy"}
		}

		if sensor.PollingStrategy == models.PollingStrategyPing && len(sensor.PollingEndpoint) == 0 {
			return &models.ValidationError{Message: "Polling endpoint is required"}
		}
	}

	return nil
}
