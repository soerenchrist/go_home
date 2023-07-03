package evaluation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/soerenchrist/go_home/internal/models"
	"github.com/soerenchrist/go_home/internal/rules"
	"github.com/soerenchrist/go_home/internal/rules/evaluation"
)

func TestDetermineUsedSensorValues_ShouldFindCorrectValues(t *testing.T) {
	expressions := []string{
		"when ${device1.sensor1.current} > 10 AND ${device2.sensor2.previous} == 1",
		"when ${device1.sensor3.current} <= 10.12 OR ${device2.sensor2.current} != true",
	}

	expected := [][]evaluation.UsedSensorValue{
		{
			{
				DeviceId: "device1",
				SensorId: "sensor1",
				Type:     evaluation.CurrentSensorValue,
			},
			{
				DeviceId: "device2",
				SensorId: "sensor2",
				Type:     evaluation.PreviousSensorValue,
			},
		},

		{
			{
				DeviceId: "device1",
				SensorId: "sensor3",
				Type:     evaluation.CurrentSensorValue,
			},
			{
				DeviceId: "device2",
				SensorId: "sensor2",
				Type:     evaluation.CurrentSensorValue,
			},
		},
	}

	for i, expression := range expressions {
		rule := &rules.Rule{When: rules.WhenExpression(expression)}

		usedSensorValues, err := evaluation.DetermineUsedSensors(rule)
		t.Logf("Used sensor values: %v", usedSensorValues)
		if err != nil {
			t.Errorf("Error while determining used sensor values: %v", err)
		}

		expectedUsedSensorValues := expected[i]

		assertUsedSensors(t, expectedUsedSensorValues, usedSensorValues)
	}
}

type FakeDatabase struct {
}

func (db FakeDatabase) ListRules() ([]rules.Rule, error) {
	return []rules.Rule{
		{
			When: rules.WhenExpression("when ${device1.sensor1.current} > 10 AND ${device2.sensor2.previous} == false"),
			Then: rules.ThenExpression("then ${device1.switch1} = true"),
			Name: "Test Rule 1",
			Id:   1,
		},
		{
			When: rules.WhenExpression("when ${device1.sensor1.previous} > 10 AND ${device2.sensor2.previous} == false"),
			Then: rules.ThenExpression("then ${device1.switch1} = true"),
			Name: "Test Rule 2",
			Id:   2,
		},
		{
			When: rules.WhenExpression("when ${device1.sensor1.previous} > 10 OR ${device2.sensor2.previous} == false"),
			Then: rules.ThenExpression("then ${device1.switch1} = true"),
			Name: "Test Rule 3",
			Id:   2,
		},
	}, nil
}

func (db FakeDatabase) AddRule(rule *rules.Rule) error {
	return nil
}

func (db FakeDatabase) GetSensor(deviceId, sensorId string) (*models.Sensor, error) {
	if deviceId == "device1" && sensorId == "sensor1" {
		return &models.Sensor{
			DeviceID: "device1",
			ID:       "sensor1",
			Type:     models.SensorTypeExternal,
			DataType: models.DataTypeInt,
			Name:     "Sensor 1",
			IsActive: true,
		}, nil
	} else if deviceId == "device2" && sensorId == "sensor2" {
		return &models.Sensor{
			DeviceID: "device2",
			ID:       "sensor2",
			Type:     models.SensorTypeExternal,
			DataType: models.DataTypeBool,
			Name:     "Sensor 2",
			IsActive: true,
		}, nil
	}
	return nil, fmt.Errorf("Sensor not found for %s.%s", deviceId, sensorId)
}

func (db FakeDatabase) GetCurrentSensorValue(deviceId, sensorId string) (*models.SensorValue, error) {
	sensorValues := map[string]*models.SensorValue{
		"device1.sensor1": {
			SensorID:  "sensor1",
			DeviceID:  "device1",
			Value:     "11",
			Timestamp: time.Now(),
		},
		"device2.sensor2": {
			SensorID:  "sensor2",
			DeviceID:  "device2",
			Value:     "true",
			Timestamp: time.Now(),
		},
	}

	key := deviceId + "." + sensorId
	if sensorValue, ok := sensorValues[key]; ok {
		return sensorValue, nil
	}

	return nil, fmt.Errorf("Sensor value not found for %s", key)
}

func (db FakeDatabase) GetPreviousSensorValue(deviceId, sensorId string) (*models.SensorValue, error) {
	sensorValues := map[string]*models.SensorValue{
		"device1.sensor1": {
			SensorID:  "sensor1",
			DeviceID:  "device1",
			Value:     "8",
			Timestamp: time.Now(),
		},
		"device2.sensor2": {
			SensorID:  "sensor2",
			DeviceID:  "device2",
			Value:     "false",
			Timestamp: time.Now(),
		},
	}

	key := deviceId + "." + sensorId
	if sensorValue, ok := sensorValues[key]; ok {
		return sensorValue, nil
	}

	return nil, fmt.Errorf("Sensor value not found for %s", key)
}

func (db FakeDatabase) GetCommand(deviceId, commandId string) (*models.Command, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (db FakeDatabase) GetDevice(deviceId string) (*models.Device, error) {
	return nil, fmt.Errorf("Not implemented")
}

func TestRuleEvaluation(t *testing.T) {
	database := FakeDatabase{}
	rulesEngine := evaluation.NewRulesEngine(database)

	rules, err := database.ListRules()
	if err != nil {
		t.Errorf("Error while listing rules: %v", err)
	}

	expectedResults := []bool{true, false, true}

	for i, rule := range rules {
		result, err := rulesEngine.EvaluateRule(&rule)

		if err != nil {
			t.Errorf("Error while evaluating rule: %v", err)
		}

		expectedResult := expectedResults[i]
		if result != expectedResult {
			t.Errorf("Expected result %v, but got %v", expectedResult, result)
		}
	}

}

func assertUsedSensors(t *testing.T, expected, got []evaluation.UsedSensorValue) {
	if len(expected) != len(got) {
		t.Errorf("Expected %d used sensor values, but got %d", len(expected), len(got))
	}

	for i, exp := range expected {
		g := got[i]

		if exp.DeviceId != g.DeviceId {
			t.Errorf("Expected device id %s, but got %s", exp.DeviceId, g.DeviceId)
		}

		if exp.SensorId != g.SensorId {
			t.Errorf("Expected sensor id %s, but got %s", exp.SensorId, g.SensorId)
		}

		if exp.Type != g.Type {
			t.Errorf("Expected sensor value type %s, but got %s", exp.Type, g.Type)
		}
	}
}
