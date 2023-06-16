package evaluation_test

import (
	"testing"

	"github.com/soerenchrist/go_home/rules"
	"github.com/soerenchrist/go_home/rules/evaluation"
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
