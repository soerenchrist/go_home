package value

import (
	"database/sql"
	"time"

	"github.com/soerenchrist/go_home/pkg/output"
)

type SensorValue struct {
	ID        uint         `json:"id"`
	SensorID  string       `json:"sensor_id"`
	DeviceID  string       `json:"device_id"`
	Value     string       `json:"value"`
	Timestamp time.Time    `json:"timestamp"`
	ExpiresAt sql.NullTime `json:"expires_at"`
}

func (sv SensorValue) ToBindingValue() output.BindingValue {
	return output.BindingValue{
		Timestamp: sv.Timestamp,
		Value:     sv.Value,
		SensorID:  sv.SensorID,
		DeviceID:  sv.DeviceID,
	}
}

type AddSensorValueRequest struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}
