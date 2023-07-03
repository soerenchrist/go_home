package rules

import (
	"fmt"
	"strings"

	"github.com/soerenchrist/go_home/internal/models"
)

type WhenExpression string
type ThenExpression string

type RulesDatabase interface {
	AddRule(rule *Rule) error
	ListRules() ([]Rule, error)
	GetSensor(deviceId, sensorId string) (*models.Sensor, error)
	GetCurrentSensorValue(deviceId, sensorId string) (*models.SensorValue, error)
	GetPreviousSensorValue(deviceId, sensorId string) (*models.SensorValue, error)
	GetCommand(deviceId, commandId string) (*models.Command, error)
	GetDevice(deviceId string) (*models.Device, error)
}

type Rule struct {
	Id   int64
	Name string
	When WhenExpression
	Then ThenExpression

	conditionAst     *Node
	actionExpression *ActionExpression
}

type Operator string
type BooleanOperator string

type ConditionExpression struct {
	SensorId string
	DeviceId string
	Variable string
	Operator Operator
	Value    string
}

type ActionExpression struct {
	DeviceId  string
	CommandId string
	Payload   string
}

type Node struct {
	Left            *Node
	Right           *Node
	BooleanOperator BooleanOperator
	Expression      *ConditionExpression
}

func (rule *Rule) ReadAction() (*ActionExpression, error) {
	if rule.actionExpression != nil {
		return rule.actionExpression, nil
	}
	tokens := strings.Split(string(rule.Then), " ")

	if len(tokens) == 0 || tokens[0] == "" {
		return nil, fmt.Errorf("invalid rule: Then Expression is empty")
	}

	if strings.ToUpper(tokens[0]) != "THEN" {
		return nil, fmt.Errorf("invalid rule: %s - Expected THEN keyword", rule.Then)
	}

	action, err := rule.parseThenExpression(tokens[1:])
	if err != nil {
		return nil, err
	}
	rule.actionExpression = action
	return action, nil
}

func (rule *Rule) ReadConditionAst() (*Node, error) {
	if rule.conditionAst != nil {
		return rule.conditionAst, nil
	}
	tokens := strings.Split(string(rule.When), " ")

	if len(tokens) == 0 || tokens[0] == "" {
		return nil, fmt.Errorf("invalid rule: When Expression is empty")
	}

	if strings.ToUpper(tokens[0]) != "WHEN" {
		return nil, fmt.Errorf("invalid rule: %s - Expected WHEN keyword", rule.When)
	}

	expression, err := rule.parseWhenExpression(tokens[1:])
	if err != nil {
		return nil, err
	}
	rule.conditionAst = expression
	return expression, nil
}

func (rule *Rule) parseWhenExpression(tokens []string) (*Node, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("invalid rule: %s - When Expression is empty", rule.When)
	}
	currentIndex := 0

	token := tokens[currentIndex]
	if !rule.isVariable(token) {
		return nil, fmt.Errorf("invalid rule: %s - Expected variable", rule.When)
	}

	deviceId, sensorId, variable, err := rule.readSensorVariable(token)
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
	sensor := &ConditionExpression{
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
	rightSide, err := rule.parseWhenExpression(tokens[currentIndex+1:])
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

func (rule *Rule) parseThenExpression(tokens []string) (*ActionExpression, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("invalid rule: %s - Then Expression is empty", rule.Then)
	}
	currentIndex := 0

	token := tokens[currentIndex]
	if !rule.isVariable(token) {
		return nil, fmt.Errorf("invalid rule: %s - Expected command variable", rule.Then)
	}

	action := &ActionExpression{}
	deviceId, commandId, err := rule.readCommandVariable(token)
	if err != nil {
		return nil, err
	}

	action.CommandId = commandId
	action.DeviceId = deviceId

	currentIndex++
	if len(tokens) == currentIndex {
		return action, nil
	}

	payload := strings.Join(tokens[currentIndex:], " ")
	action.Payload = payload

	return action, nil
}

func (rule *Rule) isVariable(token string) bool {
	return strings.HasPrefix(token, "${") && strings.HasSuffix(token, "}")
}

func (rule *Rule) readSensorVariable(token string) (deviceId string, sensorId string, variable string, err error) {
	parts := strings.Split(token[2:len(token)-1], ".")

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid variable: %s - Should consist of deviceId.sensorId.variable", token)
	}

	return parts[0], parts[1], parts[2], nil
}

func (rule *Rule) readCommandVariable(token string) (deviceId string, commandId string, err error) {
	parts := strings.Split(token[2:len(token)-1], ".")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid variable: %s - Should consist of deviceId.commandId", token)
	}

	return parts[0], parts[1], nil
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
