package rules

import (
	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
)

type RulesEngine struct {
	database db.DevicesDatabase
	sensors  map[string]*models.Sensor
}

type SensorValueType string

const (
	PreviousSensorValue SensorValueType = "previous"
	CurrentSensorValue  SensorValueType = "current"
)

type UsedSensorValue struct {
	DeviceId string
	SensorId string
	Type     SensorValueType
}

func InitRulesEngine(database db.DevicesDatabase) (*RulesEngine, error) {
	sensors, err := readSensors(database)
	if err != nil {
		return nil, err
	}
	return &RulesEngine{database: database, sensors: sensors}, nil
}

/*
*
Preload and cache all devices and sensors from the database
*/
func readSensors(database db.DevicesDatabase) (map[string]*models.Sensor, error) {
	devices, err := database.ListDevices()
	if err != nil {
		return nil, err
	}

	sensorsMap := make(map[string]*models.Sensor)
	for _, device := range devices {

		sensors, err := database.ListSensors(device.ID)
		if err != nil {
			return sensorsMap, err
		}
		for _, sensor := range sensors {
			id := sensor.DeviceID + "." + sensor.ID
			sensorsMap[id] = &sensor
		}
	}

	return sensorsMap, nil
}

func (e *RulesEngine) DetermineUsedSensors(rule *Rule) ([]UsedSensorValue, error) {
	ast, err := rule.ReadAst()
	if err != nil {
		return nil, err
	}

	usedValues := make([]UsedSensorValue, 0)
	determineUsedSensorsRec(ast, &usedValues)

	return usedValues, nil
}

func determineUsedSensorsRec(node *Node, usedValues *[]UsedSensorValue) {
	if node.Expression != nil {
		value := UsedSensorValue{node.Expression.DeviceId, node.Expression.SensorId, SensorValueType(node.Expression.Variable)}
		if !contains(*usedValues, value) {
			*usedValues = append(*usedValues, value)
		}
	}

	if node.Left != nil {
		determineUsedSensorsRec(node.Left, usedValues)
	}

	if node.Right != nil {
		determineUsedSensorsRec(node.Right, usedValues)
	}
}

func contains(values []UsedSensorValue, value UsedSensorValue) bool {
	for _, v := range values {
		if v.DeviceId == value.DeviceId && v.SensorId == value.SensorId && v.Type == value.Type {
			return true
		}
	}
	return false
}

/*
func (e *RulesEngine) CompileRule(rule *Rule, prevValues []models.SensorValue, currentValues []models.SensorValue) error {

}
*/
