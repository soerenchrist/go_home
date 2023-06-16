package rules

import (
	"fmt"
	"strings"
)

type WhenExpression string
type ThenExpression string

type Rule struct {
	Id   string
	Name string
	When WhenExpression
	Then ThenExpression
}

type Operator string

type DependentSensor struct {
	SensorId string
	DeviceId string
	Variable string
	Operator Operator
	Value    string
}

func (rule *Rule) ReadDependentSensors() ([]DependentSensor, error) {
	tokens := strings.Split(string(rule.When), " ")

	if len(tokens) == 0 || tokens[0] == "" {
		return nil, fmt.Errorf("invalid rule: %s - When Expression is empty", rule.When)
	}

	if strings.ToUpper(tokens[0]) != "WHEN" {
		return nil, fmt.Errorf("invalid rule: %s - Expected WHEN keyword", rule.When)
	}

	return rule.evaluateExpression(tokens[1:])
}

func (rule *Rule) evaluateExpression(tokens []string) ([]DependentSensor, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("invalid rule: %s - When Expression is empty", rule.When)
	}
	currentIndex := 0

	results := make([]DependentSensor, 0)
	for {
		token := tokens[currentIndex]
		if !rule.isVariable(token) {
			return nil, fmt.Errorf("invalid rule: %s - Expected variable", rule.When)
		}

		deviceId, sensorId, variable, err := rule.readVariable(token)
		if err != nil {
			return nil, err
		}

		currentIndex++
		if len(tokens) == currentIndex || !rule.isOperator(tokens[currentIndex]) {
			return nil, fmt.Errorf("invalid rule: %s - Expected operator", rule.When)
		}
		operator := Operator(tokens[currentIndex])

		currentIndex++
		if len(tokens) == currentIndex {
			return nil, fmt.Errorf("invalid rule: %s - Expected value", rule.When)
		}

		value := tokens[currentIndex]
		results = append(results, DependentSensor{
			SensorId: sensorId,
			DeviceId: deviceId,
			Variable: variable,
			Operator: operator,
			Value:    value,
		})

		currentIndex++

		if len(tokens) == currentIndex {
			break
		}

		if !rule.isBooleanOperator(tokens[currentIndex]) {
			return nil, fmt.Errorf("invalid rule: %s - Expected boolean operator", rule.When)
		}
		currentIndex++
	}
	return results, nil
}

func (rule *Rule) isVariable(token string) bool {
	return strings.HasPrefix(token, "${") && strings.HasSuffix(token, "}")
}

func (rule *Rule) readVariable(token string) (deviceId string, sensorId string, variable string, err error) {
	parts := strings.Split(token[2:len(token)-1], ".")

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid variable: %s - Should consist of deviceId.sensorId.variable", token)
	}

	return parts[0], parts[1], parts[2], nil
}

var operators = []Operator{
	Operator("=="),
	Operator("!="),
	Operator(">"),
	Operator("<"),
	Operator(">="),
	Operator("<="),
}

func (rule *Rule) isOperator(token string) bool {
	for _, operator := range operators {
		if operator == Operator(token) {
			return true
		}
	}
	return false
}

func (rule *Rule) isBooleanOperator(token string) bool {
	return strings.ToUpper(token) == "AND" || strings.ToUpper(token) == "OR"
}
