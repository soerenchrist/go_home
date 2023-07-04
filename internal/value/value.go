package value

import (
	"database/sql"
	"time"
)

type SensorValue struct {
	ID        uint         `json:"id"`
	SensorID  string       `json:"sensor_id"`
	DeviceID  string       `json:"device_id"`
	Value     string       `json:"value"`
	Timestamp time.Time    `json:"timestamp"`
	ExpiresAt sql.NullTime `json:"expires_at"`
}

type AddSensorValueRequest struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}
