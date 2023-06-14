package models

import "fmt"

type Device struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LastReached string `json:"last_reached"`
}

func (d *Device) String() string {
	return fmt.Sprintf("Device<%s %s>", d.ID, d.Name)
}

type Sensor struct {
	ID       string   `json:"id"`
	DeviceID string   `json:"device_id"`
	Name     string   `json:"name"`
	DataType DataType `json:"data_type"`
	Unit     string   `json:"unit"`
}

type DataType string

const (
	DataTypeString DataType = "string"
	DataTypeInt    DataType = "int"
	DataTypeFloat  DataType = "float"
	DataTypeBool   DataType = "bool"
)
