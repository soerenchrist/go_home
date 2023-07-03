package db

import (
	"fmt"
	"time"

	"github.com/soerenchrist/go_home/models"
)

func (db *SqliteDevicesDatabase) AddSensorValue(data *models.SensorValue) error {
	result := db.db.Create(data)
	return result.Error
}

func (db *SqliteDevicesDatabase) GetCurrentSensorValue(deviceId, sensorId string) (*models.SensorValue, error) {
	sensorVal := models.SensorValue{}
	result := db.db.Where("sensor_id = ? and device_id = ?", sensorId, deviceId).Order("timestamp desc").First(&sensorVal)

	return &sensorVal, result.Error
}

func (db *SqliteDevicesDatabase) GetPreviousSensorValue(deviceId, sensorId string) (*models.SensorValue, error) {
	sensorVals := make([]models.SensorValue, 0)
	result := db.db.Where("sensor_id = ? and device_id = ?", sensorId, deviceId).Order("timestamp desc").Limit(2).Find(&sensorVals)

	if len(sensorVals) != 2 {
		return nil, fmt.Errorf("no previous value found for sensor")
	}

	sensorVal := sensorVals[1]

	return &sensorVal, result.Error
}

func (db *SqliteDevicesDatabase) GetSensorValuesSince(deviceId, sensorId string, timestamp time.Time) ([]models.SensorValue, error) {

	values := make([]models.SensorValue, 0)
	result := db.db.Where("device_id = ? AND sensor_id = ? AND timestamp > ?", deviceId, sensorId, timestamp).Find(&values)
	return values, result.Error
}
