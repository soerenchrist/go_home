package db

import (
	"github.com/soerenchrist/go_home/internal/command"
	"github.com/soerenchrist/go_home/internal/errors"
)

func (db *SqliteDevicesDatabase) GetCommand(deviceId, commandId string) (*command.Command, error) {
	command := command.Command{}
	result := db.db.Where("id = ? and device_id = ?", commandId, deviceId).First(&command)
	return &command, result.Error
}

func (db *SqliteDevicesDatabase) ListCommands(deviceId string) ([]command.Command, error) {
	commands := make([]command.Command, 0)
	result := db.db.Where("device_id = ?", deviceId).Find(&commands)
	return commands, result.Error
}

func (db *SqliteDevicesDatabase) AddCommand(command *command.Command) error {
	result := db.db.Create(command)
	return result.Error
}

func (db *SqliteDevicesDatabase) DeleteCommand(deviceId, commandId string) error {
	result := db.db.Where("id = ? and device_id = ?", commandId, deviceId).Delete(&command.Command{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &errors.NotFoundError{Message: "Command not found"}
	}
	return result.Error
}
