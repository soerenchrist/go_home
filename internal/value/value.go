package value

import "time"

type SensorValue struct {
	SensorID  string    `json:"sensor_id"`
	DeviceID  string    `json:"device_id"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type AddSensorValueRequest struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}
