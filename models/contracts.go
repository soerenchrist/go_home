package models

type CreateDeviceRequest struct {
	Name string `json:"name"`
}

type CreateSensorRequest struct {
	Name     string   `json:"name"`
	DataType DataType `json:"data_type"`
	Unit     string   `json:"unit"`
}

type AddSensorValueRequest struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}
