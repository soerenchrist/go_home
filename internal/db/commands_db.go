package db

import (
	"log"

	"github.com/soerenchrist/go_home/internal/models"
)

func (db *SqliteDevicesDatabase) GetCommand(deviceId, commandId string) (*models.Command, error) {
	command := models.Command{}
	result := db.db.Where("id = ? and device_id = ?", commandId, deviceId).First(&command)
	log.Println(result.Error)

	return &command, result.Error
}

func (db *SqliteDevicesDatabase) ListCommands(deviceId string) ([]models.Command, error) {
	commands := make([]models.Command, 0)
	result := db.db.Where("device_id = ?", deviceId).Find(&commands)
	return commands, result.Error
}

func (db *SqliteDevicesDatabase) AddCommand(command *models.Command) error {
	result := db.db.Create(command)
	return result.Error
}

func (db *SqliteDevicesDatabase) DeleteCommand(deviceId, commandId string) error {
	result := db.db.Where("id = ? and device_id = ?", commandId, deviceId).Delete(&models.Command{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &models.NotFoundError{Message: "Command not found"}
	}
	return result.Error
}
