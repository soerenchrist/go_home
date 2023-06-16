package rules

import (
	"github.com/soerenchrist/go_home/models"
)

type RulesEngine struct {
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
