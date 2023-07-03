package db

import (
	"github.com/soerenchrist/go_home/models"
)

func (db *SqliteDevicesDatabase) GetSensor(deviceId, sensorId string) (*models.Sensor, error) {
	sensor := models.Sensor{}
	result := db.db.Where("id = ? and device_id = ?", sensorId, deviceId).First(&sensor)

	return &sensor, result.Error
}

func (db *SqliteDevicesDatabase) ListSensors(deviceId string) ([]models.Sensor, error) {
	sensors := make([]models.Sensor, 0)
	result := db.db.Where("device_id = ?", deviceId).Find(&sensors)
	return sensors, result.Error
}

func (db *SqliteDevicesDatabase) ListPollingSensors() ([]models.Sensor, error) {
	sensors := make([]models.Sensor, 0)
	result := db.db.Where("type = 'polling'").Find(&sensors)
	return sensors, result.Error
}

func (db *SqliteDevicesDatabase) AddSensor(sensor *models.Sensor) error {
	result := db.db.Create(sensor)
	return result.Error
}

func (db *SqliteDevicesDatabase) DeleteSensor(deviceId, sensorId string) error {
	result := db.db.Where("id = ? and device_id = ?", sensorId, deviceId).Delete(&models.Sensor{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &models.NotFoundError{Message: "Sensor not found"}
	}
	return nil
}
