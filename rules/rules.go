package rules

import (
	"fmt"
	"strings"
)

type WhenExpression string
type ThenExpression string

type Rule struct {
	Id   int64
	Name string
	When WhenExpression
	Then ThenExpression
}

type Operator string
type BooleanOperator string

type Expression struct {
	SensorId string
	DeviceId string
	Variable string
	Operator Operator
	Value    string
}

type Node struct {
	Left            *Node
	Right           *Node
	BooleanOperator BooleanOperator
	Expression      *Expression
}

func (rule *Rule) ReadAst() (*Node, error) {
	tokens := strings.Split(string(rule.When), " ")

	if len(tokens) == 0 || tokens[0] == "" {
		return nil, fmt.Errorf("invalid rule: When Expression is empty")
	}

	if strings.ToUpper(tokens[0]) != "WHEN" {
		return nil, fmt.Errorf("invalid rule: %s - Expected WHEN keyword", rule.When)
	}

	expression, err := rule.evaluateExpression(tokens[1:])
	return expression, err
}

func (rule *Rule) evaluateExpression(tokens []string) (*Node, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("invalid rule: %s - When Expression is empty", rule.When)
	}
	currentIndex := 0

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
	sensor := &Expression{
		SensorId: sensorId,
		DeviceId: deviceId,
		Variable: variable,
		Operator: operator,
		Value:    value,
	}

	currentIndex++

	// Expression is finished
	if len(tokens) == currentIndex {
		return &Node{
			Expression: sensor,
		}, nil
	}

	// Expression is not finished
	if !rule.isBooleanOperator(tokens[currentIndex]) {
		return nil, fmt.Errorf("invalid rule: %s - Expected boolean operator", rule.When)
	}

	boolOp := BooleanOperator(tokens[currentIndex])
	rightSide, err := rule.evaluateExpression(tokens[currentIndex+1:])
	if err != nil {
		return nil, err
	}
	return &Node{
		Left: &Node{
			Expression: sensor,
		},
		BooleanOperator: boolOp,
		Right:           rightSide,
	}, nil
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
