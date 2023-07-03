package db

import (
	"github.com/soerenchrist/go_home/internal/errors"
	"github.com/soerenchrist/go_home/internal/sensor"
)

func (db *SqliteDevicesDatabase) GetSensor(deviceId, sensorId string) (*sensor.Sensor, error) {
	sensor := sensor.Sensor{}
	result := db.db.Where("id = ? and device_id = ?", sensorId, deviceId).First(&sensor)

	return &sensor, result.Error
}

func (db *SqliteDevicesDatabase) ListSensors(deviceId string) ([]sensor.Sensor, error) {
	sensors := make([]sensor.Sensor, 0)
	result := db.db.Where("device_id = ?", deviceId).Find(&sensors)
	return sensors, result.Error
}

func (db *SqliteDevicesDatabase) ListPollingSensors() ([]sensor.Sensor, error) {
	sensors := make([]sensor.Sensor, 0)
	result := db.db.Where("type = 'polling'").Find(&sensors)
	return sensors, result.Error
}

func (db *SqliteDevicesDatabase) AddSensor(sensor *sensor.Sensor) error {
	result := db.db.Create(sensor)
	return result.Error
}

func (db *SqliteDevicesDatabase) DeleteSensor(deviceId, sensorId string) error {
	result := db.db.Where("id = ? and device_id = ?", sensorId, deviceId).Delete(&sensor.Sensor{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &errors.NotFoundError{Message: "Sensor not found"}
	}
	return nil
}
