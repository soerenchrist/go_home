package sensor

import "time"

type Sensor struct {
	ID              string          `json:"id" gorm:"primaryKey"`
	DeviceID        string          `json:"device_id" gorm:"primaryKey"`
	Name            string          `json:"name"`
	DataType        DataType        `json:"data_type"`
	Unit            string          `json:"unit"`
	IsActive        bool            `json:"is_active"`
	Type            SensorType      `json:"type"`
	PollingInterval int             `json:"polling_interval"`
	PollingEndpoint string          `json:"polling_endpoint"`
	PollingStrategy PollingStrategy `json:"polling_strategy"`

	CreatedAt time.Time
	UpdatedAt time.Time
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

type CreateSensorRequest struct {
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	DataType        DataType        `json:"data_type"`
	Unit            string          `json:"unit"`
	Type            SensorType      `json:"type"`
	PollingInterval int             `json:"polling_interval"`
	PollingEndpoint string          `json:"polling_endpoint"`
	PollingStrategy PollingStrategy `json:"polling_strategy"`
}
