package db

import (
	"time"

	"github.com/soerenchrist/go_home/internal/command"
	"github.com/soerenchrist/go_home/internal/device"
	"github.com/soerenchrist/go_home/internal/rules"
	"github.com/soerenchrist/go_home/internal/sensor"
	"github.com/soerenchrist/go_home/internal/value"
	"gorm.io/gorm"
)

type Database interface {
	AddDevice(entity *device.Device) error
	GetDevice(id string) (*device.Device, error)
	DeleteDevice(id string) error
	ListDevices() ([]device.Device, error)
	ListSensors(deviceId string) ([]sensor.Sensor, error)
	AddSensor(sensor *sensor.Sensor) error
	GetSensor(deviceId, sensorId string) (*sensor.Sensor, error)
	DeleteSensor(deviceId, sensorId string) error

	ListPollingSensors() ([]sensor.Sensor, error)

	AddSensorValue(sensorValue *value.SensorValue) error
	GetSensorValuesSince(deviceId, sensorId string, timestamp time.Time) ([]value.SensorValue, error)
	GetCurrentSensorValue(deviceId, sensorId string) (*value.SensorValue, error)
	GetPreviousSensorValue(deviceId, sensorId string) (*value.SensorValue, error)

	AddCommand(command *command.Command) error
	GetCommand(deviceId, commandId string) (*command.Command, error)
	ListCommands(deviceId string) ([]command.Command, error)
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
	db.db.AutoMigrate(&command.Command{}, &device.Device{}, &sensor.Sensor{}, &value.SensorValue{}, &rules.Rule{})
	return nil
}

func (database *SqliteDevicesDatabase) SeedDatabase() {
	device1 := &device.Device{ID: "1", Name: "My Device 1"}
	sensor1 := &sensor.Sensor{ID: "S1", Name: "Temperature", DeviceID: "1", DataType: sensor.DataTypeFloat, Type: sensor.SensorTypeExternal, IsActive: true, Unit: "Celsius", PollingInterval: 0, RetainmentPeriodSeconds: 3600}
	sensor2 := &sensor.Sensor{ID: "S2", Name: "Availability", DeviceID: "1", DataType: sensor.DataTypeBool, Type: sensor.SensorTypePolling, IsActive: true, Unit: "", PollingInterval: 120, PollingEndpoint: "localhost", PollingStrategy: "ping"}
	template := `{"device": "{{.device_id}}", "command": "{{.command_id}}", "payload": "{{.p_payload}}"}`
	command1 := &command.Command{ID: "C1", Name: "Turn on", DeviceID: "1", PayloadTemplate: template, Endpoint: "http://localhost:8080/echo", Method: "POST"}

	device2 := &device.Device{ID: "2", Name: "My Device 2"}
	sensor3 := &sensor.Sensor{ID: "S3", Name: "Filling Level", DeviceID: "2", DataType: sensor.DataTypeInt, Type: sensor.SensorTypeExternal, IsActive: true, Unit: "%", PollingInterval: 0}

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
