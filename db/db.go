package db

import "github.com/soerenchrist/mini_home/models"

type DevicesDatabase interface {
	AddDevice(entity models.Device) error
	GetDevice(id string) (models.Device, error)
	DeleteDevice(id string) error
	ListDevices() ([]models.Device, error)
	ListSensors(deviceId string) ([]models.Sensor, error)
	AddSensor(sensor models.Sensor) error
	Close() error
}
