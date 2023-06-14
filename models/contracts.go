package models

type CreateDeviceRequest struct {
	Name string `json:"name"`
}

type CreateSensorRequest struct {
	Name            string          `json:"name"`
	DataType        DataType        `json:"data_type"`
	Unit            string          `json:"unit"`
	Type            SensorType      `json:"type"`
	PollingInterval int             `json:"polling_interval"`
	PollingEndpoint string          `json:"polling_endpoint"`
	PollingStrategy PollingStrategy `json:"polling_strategy"`
}

type AddSensorValueRequest struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

type CreateCommandRequest struct {
	Name            string `json:"name"`
	PayloadTemplate string `json:"payload_template"`
	Endpoint        string `json:"endpoint"`
	Method          string `json:"method"`
}
