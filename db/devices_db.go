package db

import (
	"github.com/soerenchrist/go_home/models"
)

func (db *SqliteDevicesDatabase) AddDevice(device *models.Device) error {
	result := db.db.Create(device)
	return result.Error
}

func (db *SqliteDevicesDatabase) GetDevice(id string) (*models.Device, error) {
	device := models.Device{}
	result := db.db.First(&device, id)
	return &device, result.Error
}

func (db *SqliteDevicesDatabase) DeleteDevice(id string) error {
	result := db.db.Delete(&models.Device{ID: id})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &models.NotFoundError{Message: "Device not found"}
	}
	return nil
}

func (db *SqliteDevicesDatabase) ListDevices() ([]models.Device, error) {
	devices := make([]models.Device, 0)
	result := db.db.Find(&devices)
	return devices, result.Error
}
