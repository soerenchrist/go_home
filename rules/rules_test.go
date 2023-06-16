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

		_, err := rule.ReadAst()

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
		"when ${1.S2.curr} != true OR ${1.S2.prev} == false",
	}

	expectedResult := []rules.Node{
		{
			Expression: &rules.Expression{
				SensorId: "S1",
				DeviceId: "1",
				Variable: "curr",
				Operator: ">",
				Value:    "1",
			},
		},
		{
			Left: &rules.Node{
				Expression: &rules.Expression{
					SensorId: "S1",
					DeviceId: "1",
					Variable: "curr",
					Operator: ">",
					Value:    "1",
				},
			},
			BooleanOperator: "AND",
			Right: &rules.Node{
				Expression: &rules.Expression{
					SensorId: "S1",
					DeviceId: "1",
					Variable: "prev",
					Operator: "<",
					Value:    "2",
				},
			},
		},
		{
			Left: &rules.Node{
				Expression: &rules.Expression{
					SensorId: "S2",
					DeviceId: "1",
					Variable: "curr",
					Operator: "!=",
					Value:    "true",
				},
			},
			BooleanOperator: "OR",
			Right: &rules.Node{
				Expression: &rules.Expression{
					SensorId: "S2",
					DeviceId: "1",
					Variable: "prev",
					Operator: "==",
					Value:    "false",
				},
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

		result, err := rule.ReadAst()

		if err != nil {
			t.Errorf("Expression %s, got error, expected none", expression)
			return
		}

		expectedResult := expectedResult[i]

		assertResult(t, &expectedResult, result)
	}
}

func assertResult(t *testing.T, expected, got *rules.Node) {
	// assert current node
	assertSensor(t, expected.Expression, got.Expression)
	if expected.BooleanOperator != got.BooleanOperator {
		t.Errorf("Expected boolean operator '%s', but got '%s'", expected.BooleanOperator, got.BooleanOperator)
	}

	if expected.Left == nil && got.Left != nil {
		t.Errorf("Expected left node to be nil, but got '%v'", got.Left)
	}

	// assert left node
	if expected.Left != nil {
		assertResult(t, expected.Left, got.Left)
	}

	if expected.Right == nil && got.Right != nil {
		t.Errorf("Expected right node to be nil, but got '%v'", got.Right)
	}

	// assert right node
	if expected.Right != nil {
		assertResult(t, expected.Right, got.Right)
	}
}

func assertSensor(t *testing.T, expected, got *rules.Expression) {
	if expected == nil && got == nil {
		return
	}
	if (expected == nil) != (got == nil) {
		t.Errorf("Expected sensor '%v', but got '%v'", expected, got)
		return
	}
	if expected.SensorId != got.SensorId {
		t.Errorf("Expected sensorId '%s', but got '%s'", expected.SensorId, got.SensorId)
	}

	if expected.DeviceId != got.DeviceId {
		t.Errorf("Expected deviceId '%s', but got '%s'", expected.DeviceId, got.DeviceId)
	}

	if expected.Variable != got.Variable {
		t.Errorf("Expected variable '%s', but got '%s'", expected.Variable, got.Variable)
	}

	if expected.Operator != got.Operator {
		t.Errorf("Expected operator '%s', but got '%s'", expected.Operator, got.Operator)
	}

	if expected.Value != got.Value {
		t.Errorf("Expected value '%s', but got '%s'", expected.Value, got.Value)
	}
}
