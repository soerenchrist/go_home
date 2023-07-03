package db

import (
	"time"

	"github.com/soerenchrist/go_home/internal/models"
	"github.com/soerenchrist/go_home/internal/rules"
	"gorm.io/gorm"
)

type Database interface {
	AddDevice(entity *models.Device) error
	GetDevice(id string) (*models.Device, error)
	DeleteDevice(id string) error
	ListDevices() ([]models.Device, error)
	ListSensors(deviceId string) ([]models.Sensor, error)
	AddSensor(sensor *models.Sensor) error
	GetSensor(deviceId, sensorId string) (*models.Sensor, error)
	DeleteSensor(deviceId, sensorId string) error

	ListPollingSensors() ([]models.Sensor, error)

	AddSensorValue(sensorValue *models.SensorValue) error
	GetSensorValuesSince(deviceId, sensorId string, timestamp time.Time) ([]models.SensorValue, error)
	GetCurrentSensorValue(deviceId, sensorId string) (*models.SensorValue, error)
	GetPreviousSensorValue(deviceId, sensorId string) (*models.SensorValue, error)

	AddCommand(command *models.Command) error
	GetCommand(deviceId, commandId string) (*models.Command, error)
	ListCommands(deviceId string) ([]models.Command, error)
	DeleteCommand(deviceId, commandId string) error

	ListRules() ([]rules.Rule, error)
	AddRule(rule *rules.Rule) error

	SeedDatabase()
}

type SqliteDevicesDatabase struct {
	db *gorm.DB
}

func NewDevicesDatabase(db *gorm.DB) (*SqliteDevicesDatabase, error) {
	database := &SqliteDevicesDatabase{db: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (db *SqliteDevicesDatabase) createTables() error {
	db.db.AutoMigrate(&models.Command{}, &models.Device{}, &models.Sensor{}, &models.SensorValue{}, &rules.Rule{})
	return nil
}

func (database *SqliteDevicesDatabase) SeedDatabase() {
	device1 := &models.Device{ID: "1", Name: "My Device 1"}
	sensor1 := &models.Sensor{ID: "S1", Name: "Temperature", DeviceID: "1", DataType: models.DataTypeFloat, Type: models.SensorTypeExternal, IsActive: true, Unit: "Celsius", PollingInterval: 0}
	sensor2 := &models.Sensor{ID: "S2", Name: "Availability", DeviceID: "1", DataType: models.DataTypeBool, Type: models.SensorTypePolling, IsActive: true, Unit: "", PollingInterval: 10, PollingEndpoint: "localhost", PollingStrategy: "ping"}
	template := `{"device": "{{.device_id}}", "command": "{{.command_id}}", "payload": "{{.p_payload}}"}`
	command1 := &models.Command{ID: "C1", Name: "Turn on", DeviceID: "1", PayloadTemplate: template, Endpoint: "http://localhost:8080/echo", Method: "POST"}

	device2 := &models.Device{ID: "2", Name: "My Device 2"}
	sensor3 := &models.Sensor{ID: "S3", Name: "Filling Level", DeviceID: "2", DataType: models.DataTypeInt, Type: models.SensorTypeExternal, IsActive: true, Unit: "%", PollingInterval: 0}

	devices, err := database.ListDevices()
	if err != nil {
		panic(err)
	}
	if len(devices) > 0 {
		return
	}

	if err := database.AddDevice(device1); err != nil {
		panic(err)
	}
	if err := database.AddDevice(device2); err != nil {
		panic(err)
	}

	if err := database.AddSensor(sensor1); err != nil {
		panic(err)
	}
	if err := database.AddSensor(sensor2); err != nil {
		panic(err)
	}
	if err := database.AddSensor(sensor3); err != nil {
		panic(err)
	}
	if err := database.AddCommand(command1); err != nil {
		panic(err)
	}

	rule := &rules.Rule{
		Name: "Turn on light when temperature is below 20",
		When: rules.WhenExpression("when ${1.S1.current} < 20 AND ${1.S1.previous} >= 20"),
		Then: rules.ThenExpression("then ${1.C1} {\"p_payload\": \"on\"}"),
	}

	if err := database.AddRule(rule); err != nil {
		panic(err)
	}
}