package rules_test

import (
	"strings"
	"testing"

	"github.com/soerenchrist/go_home/rules"
)

func TestInvalidWhenExpressions(t *testing.T) {
	expressions := []string{
		"",
		"Hello world",
		"when",
		"when 1",
		"when ${sensor}",
		"when ${}",
		"when ${1.S1.curr}",
		"when ${1.S1.curr} ! 1",
		"when ${1.S1.curr} >",
		"when ${1.S1.curr} > 1 something",
		"when ${1.S1.curr} > 1 AND ${1.S1.curr}",
	}

	expectedMessages := []string{
		"When Expression is empty",
		"Expected WHEN keyword",
		"When Expression is empty",
		"Expected variable",
		"Should consist of deviceId.sensorId.variable",
		"Should consist of deviceId.sensorId.variable",
		"Expected operator",
		"Expected operator",
		"Expected value",
		"Expected boolean operator",
		"Expected operator",
	}

	for i, expression := range expressions {
		rule := &rules.Rule{
			Id:   "1",
			Name: "Test Rule",
			When: rules.WhenExpression(expression),
			Then: rules.ThenExpression(""),
		}

		_, err := rule.ReadDependentSensors()

		if err == nil {
			t.Errorf("Expression %s should give error, but got none", expression)
			return
		}

		expectedMessage := expectedMessages[i]

		if !strings.HasSuffix(err.Error(), expectedMessage) {
			t.Errorf("Expected '%s', but got '%s'", expectedMessage, err.Error())
		}
	}
}

func TestReadDependentSensors_ShouldReturnCorrectValues(t *testing.T) {
	expressions := []string{
		"when ${1.S1.curr} > 1",
		"when ${1.S1.curr} > 1 AND ${1.S1.prev} < 2",
		"when ${1.S2.curr} != true AND ${1.S2.prev} == false",
	}

	expectedResult := [][]rules.DependentSensor{
		{
			{
				SensorId: "S1",
				DeviceId: "1",
				Variable: "curr",
				Operator: rules.Operator(">"),
				Value:    "1",
			},
		},
		{
			{
				SensorId: "S1",
				DeviceId: "1",
				Variable: "curr",
				Operator: rules.Operator(">"),
				Value:    "1",
			},
			{
				SensorId: "S1",
				DeviceId: "1",
				Variable: "prev",
				Operator: rules.Operator("<"),
				Value:    "2",
			},
		},
		{
			{
				SensorId: "S2",
				DeviceId: "1",
				Variable: "curr",
				Operator: rules.Operator("!="),
				Value:    "true",
			},
			{
				SensorId: "S2",
				DeviceId: "1",
				Variable: "prev",
				Operator: rules.Operator("=="),
				Value:    "false",
			},
		},
	}

	for i, expression := range expressions {
		rule := &rules.Rule{
			Id:   "1",
			Name: "Test Rule",
			When: rules.WhenExpression(expression),
			Then: rules.ThenExpression(""),
		}

		result, err := rule.ReadDependentSensors()

		if err != nil {
			t.Errorf("Expression %s, got error, expected none", expression)
			return
		}

		expectedResult := expectedResult[i]

		assertResult(t, expectedResult, result)
	}
}

func assertResult(t *testing.T, expected []rules.DependentSensor, got []rules.DependentSensor) {
	if len(expected) != len(got) {
		t.Errorf("Expected %d results, but got %d", len(expected), len(got))
		return
	}

	for i, expectedSensor := range expected {
		if expectedSensor.SensorId != got[i].SensorId {
			t.Errorf("Expected sensorId '%s', but got '%s'", expectedSensor.SensorId, got[i].SensorId)
		}

		if expectedSensor.DeviceId != got[i].DeviceId {
			t.Errorf("Expected deviceId '%s', but got '%s'", expectedSensor.DeviceId, got[i].DeviceId)
		}

		if expectedSensor.Variable != got[i].Variable {
			t.Errorf("Expected variable '%s', but got '%s'", expectedSensor.Variable, got[i].Variable)
		}

		if expectedSensor.Operator != got[i].Operator {
			t.Errorf("Expected operator '%s', but got '%s'", expectedSensor.Operator, got[i].Operator)
		}

		if expectedSensor.Value != got[i].Value {
			t.Errorf("Expected value '%s', but got '%s'", expectedSensor.Value, got[i].Value)
		}
	}
}
