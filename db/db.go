package db

import "github.com/soerenchrist/mini_home/models"

type DevicesDatabase interface {
	AddDevice(entity models.Device) error
	GetDevice(id string) (models.Device, error)
	DeleteDevice(id string) error
	ListDevices() ([]models.Device, error)
	ListSensors(deviceId string) ([]models.Sensor, error)
	AddSensor(sensor models.Sensor) error
	GetSensor(deviceId, sensorId string) (models.Sensor, error)
	DeleteSensor(deviceId, sensorId string) error

	ListPollingSensors() ([]models.Sensor, error)

	AddSensorValue(sensorValue models.SensorValue) error
	Close() error
	SeedDatabase()
}
