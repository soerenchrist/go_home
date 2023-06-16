package evaluation

import (
	"log"

	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
	"github.com/soerenchrist/go_home/rules"
)

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

type RulesEngine struct {
	lookupTable map[string][]rules.Rule
}

func NewRulesEngine(database db.DevicesDatabase) *RulesEngine {
	lookupTable, err := buildLookupTable(database)
	if err != nil {
		panic(err)
	}
	return &RulesEngine{lookupTable: lookupTable}
}

func (engine *RulesEngine) ListenForValues(sensorsChannel chan models.Sensor) {
	log.Println("Listening for sensor values...")
	for {
		sensor := <-sensorsChannel
		log.Printf("Received sensor value: %v\n", sensor)
	}
}

func buildLookupTable(database db.DevicesDatabase) (map[string][]rules.Rule, error) {
	lookupTable := make(map[string][]rules.Rule)

	allRules, err := database.ListRules()
	if err != nil {
		return nil, err
	}

	for _, rule := range allRules {
		usedSensors, err := DetermineUsedSensors(&rule)
		if err != nil {
			return nil, err
		}

		for _, usedSensor := range usedSensors {
			key := usedSensor.DeviceId + "." + usedSensor.SensorId
			if _, ok := lookupTable[key]; !ok {
				lookupTable[key] = make([]rules.Rule, 0)
			}
			lookupTable[key] = append(lookupTable[key], rule)
		}
	}

	return lookupTable, nil
}

func DetermineUsedSensors(rule *rules.Rule) ([]UsedSensorValue, error) {
	ast, err := rule.ReadAst()
	if err != nil {
		return nil, err
	}

	usedValues := make([]UsedSensorValue, 0)
	determineUsedSensorsRec(ast, &usedValues)

	return usedValues, nil
}

func determineUsedSensorsRec(node *rules.Node, usedValues *[]UsedSensorValue) {
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
