package db

import "github.com/soerenchrist/mini_home/models"

type DevicesDatabase interface {
	Add(entity models.Device) error
	List() ([]models.Device, error)
	ListSensors(deviceId string) ([]models.Sensor, error)
	Close() error
}
