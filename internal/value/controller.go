package value

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/go_home/internal/device"
	"github.com/soerenchrist/go_home/internal/errors"
	"github.com/soerenchrist/go_home/internal/sensor"
	"github.com/soerenchrist/go_home/internal/util"
)

type SensorValuesDatabase interface {
	GetSensorValuesSince(deviceId string, sensorId string, since time.Time) ([]SensorValue, error)
	GetSensor(deviceId string, sensorId string) (*sensor.Sensor, error)
	GetDevice(deviceId string) (*device.Device, error)
	GetCurrentSensorValue(deviceId string, sensorId string) (*SensorValue, error)
	AddSensorValue(sensorValue *SensorValue) error
}

type SensorValuesController struct {
	database       SensorValuesDatabase
	outputBindings *OutputBindings
}

func NewController(database SensorValuesDatabase, outputBindings *OutputBindings) *SensorValuesController {
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

	var request AddSensorValueRequest

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

	timestamp, _ := time.Parse(time.RFC3339, request.Timestamp)
	var expiry sql.NullTime

	if sensor.RetainmentPeriodSeconds > 0 {
		expiry = sql.NullTime{
			Time:  timestamp.Add(time.Duration(sensor.RetainmentPeriodSeconds) * time.Second),
			Valid: true,
		}
	}
	sensorValue := &SensorValue{
		Value:     request.Value,
		Timestamp: timestamp,
		DeviceID:  device.ID,
		SensorID:  sensor.ID,
		ExpiresAt: expiry,
	}

	err = c.database.AddSensorValue(sensorValue)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.outputBindings.Push(*sensorValue)
	context.JSON(201, sensorValue)
}

func (c *SensorValuesController) getSensorAndDevice(context *gin.Context) (*sensor.Sensor, *device.Device, error) {
	deviceId := context.Param("deviceId")
	sensorId := context.Param("sensorId")

	var d *device.Device
	var s *sensor.Sensor

	var err error

	if d, err = c.database.GetDevice(deviceId); err != nil {
		return &sensor.Sensor{}, &device.Device{}, &errors.NotFoundError{Message: "Device not found"}
	}

	if s, err = c.database.GetSensor(deviceId, sensorId); err != nil {
		return &sensor.Sensor{}, &device.Device{}, &errors.NotFoundError{Message: "Sensor not found"}
	}

	return s, d, nil
}

func (c *SensorValuesController) validateSensorData(s *sensor.Sensor, request *AddSensorValueRequest) error {
	if s.DataType == sensor.DataTypeInt {
		if _, err := strconv.Atoi(request.Value); err != nil {
			return &errors.ValidationError{Message: "Sensor value is not an int"}
		}
	} else if s.DataType == sensor.DataTypeFloat {
		if _, err := strconv.ParseFloat(request.Value, 64); err != nil {
			return &errors.ValidationError{Message: "Sensor value is not a float"}
		}
	} else if s.DataType == sensor.DataTypeBool {
		if _, err := strconv.ParseBool(request.Value); err != nil {
			return &errors.ValidationError{Message: "Sensor value is not a bool"}
		}
	}

	if s.Type == sensor.SensorTypePolling {
		return &errors.ValidationError{Message: "Sending values to a polling sensor is not allowed"}
	}

	if !s.IsActive {
		return &errors.ValidationError{Message: "Sensor is not active"}
	}

	if len(request.Timestamp) > 0 {
		if err := util.ValidateTimestamp(request.Timestamp); err != nil {
			return &errors.ValidationError{Message: fmt.Sprintf("%s is not a valid RFC3339 timestamp", request.Timestamp)}
		}
	}

	return nil
}
