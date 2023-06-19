package models

import "fmt"

type Device struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (d *Device) String() string {
	return fmt.Sprintf("Device<%s %s>", d.ID, d.Name)
}

type Sensor struct {
	ID              string          `json:"id"`
	DeviceID        string          `json:"device_id"`
	Name            string          `json:"name"`
	DataType        DataType        `json:"data_type"`
	Unit            string          `json:"unit"`
	IsActive        bool            `json:"is_active"`
	Type            SensorType      `json:"type"`
	PollingInterval int             `json:"polling_interval"`
	PollingEndpoint string          `json:"polling_endpoint"`
	PollingStrategy PollingStrategy `json:"polling_strategy"`
}

type PollingStrategy string

const (
	PollingStrategyPing PollingStrategy = "ping"
)

type DataType string

const (
	DataTypeString DataType = "string"
	DataTypeInt    DataType = "int"
	DataTypeFloat  DataType = "float"
	DataTypeBool   DataType = "bool"
)

type SensorType string

const (
	SensorTypeExternal SensorType = "external"
	SensorTypePolling  SensorType = "polling"
)

type SensorValue struct {
	SensorID  string `json:"sensor_id"`
	DeviceID  string `json:"device_id"`
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}
