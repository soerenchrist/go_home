package db

import (
	"github.com/soerenchrist/go_home/internal/device"
	"github.com/soerenchrist/go_home/internal/errors"
)

func (db *SqliteDevicesDatabase) AddDevice(device *device.Device) error {
	result := db.db.Create(device)
	return result.Error
}

func (db *SqliteDevicesDatabase) GetDevice(id string) (*device.Device, error) {
	device := device.Device{}
	result := db.db.First(&device, id)
	return &device, result.Error
}

func (db *SqliteDevicesDatabase) DeleteDevice(id string) error {
	result := db.db.Delete(&device.Device{ID: id})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &errors.NotFoundError{Message: "Device not found"}
	}
	return nil
}

func (db *SqliteDevicesDatabase) ListDevices() ([]device.Device, error) {
	devices := make([]device.Device, 0)
	result := db.db.Find(&devices)
	return devices, result.Error
}
