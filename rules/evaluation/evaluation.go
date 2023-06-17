package evaluation

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/soerenchrist/go_home/models"
	"github.com/soerenchrist/go_home/rules"
)

type SensorValueType string

const (
	PreviousSensorValue SensorValueType = "previous"
	CurrentSensorValue  SensorValueType = "current"
)

type RulesDatabase interface {
	ListRules() ([]rules.Rule, error)
	GetSensor(deviceId, sensorId string) (*models.Sensor, error)
	GetCurrentSensorValue(deviceId, sensorId string) (*models.SensorValue, error)
	GetPreviousSensorValue(deviceId, sensorId string) (*models.SensorValue, error)
	GetCommand(deviceId, commandId string) (*models.Command, error)
	GetDevice(deviceId string) (*models.Device, error)
}

type UsedSensorValue struct {
	DeviceId string
	SensorId string
	Type     SensorValueType
}

type RulesEngine struct {
	database    RulesDatabase
	lookupTable map[string][]rules.Rule
}

func NewRulesEngine(database RulesDatabase) *RulesEngine {
	lookupTable, err := buildLookupTable(database)
	if err != nil {
		panic(err)
	}
	return &RulesEngine{lookupTable: lookupTable, database: database}
}

func (engine *RulesEngine) ListenForValues(sensorsChannel chan models.SensorValue) {
	log.Println("Listening for sensor values...")
	for {
		sensor := <-sensorsChannel

		key := sensor.DeviceID + "." + sensor.SensorID
		if rules, ok := engine.lookupTable[key]; ok {
			for _, rule := range rules {
				log.Printf("Evaluating rule %v\n", rule)
				evalResult, err := engine.EvaluateRule(&rule)
				if err != nil {
					log.Printf("Error evaluating rule: %v\n", err)
					continue
				}
				log.Printf("Rule '%s' evaluated to %v\n", rule.Name, evalResult)
				if evalResult {
					err := engine.executeRule(rule)
					if err != nil {
						log.Printf("Error executing rule: %v\n", err)
					}
				}
			}
		}
	}
}

func (engine *RulesEngine) executeRule(rule rules.Rule) error {
	action, err := rule.ReadAction()
	if err != nil {
		return fmt.Errorf("error reading action: %v", err)
	}

	device, err := engine.database.GetDevice(action.DeviceId)
	if err != nil {
		return fmt.Errorf("error reading device: %v", err)
	}

	command, err := engine.database.GetCommand(action.DeviceId, action.CommandId)
	if err != nil {
		return fmt.Errorf("error reading command: %v", err)
	}

	log.Printf("Executing command: %v\n", command)

	var params models.CommandParameters
	if action.Payload != "" {
		err := json.Unmarshal([]byte(action.Payload), &params)
		if err != nil {
			return err
		}
	}

	resp, err := command.Invoke(device, &params)
	if err != nil {
		return fmt.Errorf("error invoking command: %v", err)
	}

	log.Printf("Command response status: %d \n", resp.StatusCode)
	return nil
}

func (engine *RulesEngine) EvaluateRule(rule *rules.Rule) (bool, error) {
	deps, err := DetermineUsedSensors(rule)
	if err != nil {
		return false, err
	}

	values, err := engine.readDependentValues(deps)
	if err != nil {
		return false, err
	}

	for key, value := range values {
		log.Printf("Sensor value: %s = %v\n", key, value)
	}

	ast, err := rule.ReadConditionAst()
	if err != nil {
		return false, err
	}

	return engine.evaluateAst(ast, values)
}

func (engine *RulesEngine) evaluateAst(ast *rules.Node, values map[string]string) (bool, error) {
	if ast.Expression != nil {
		return engine.evaluateExpression(ast.Expression, values)
	}

	var leftVal, rightVal bool
	var err error
	if ast.Left != nil {
		leftVal, err = engine.evaluateAst(ast.Left, values)
		if err != nil {
			return false, err
		}
	}

	if ast.Right != nil {
		rightVal, err = engine.evaluateAst(ast.Right, values)
		if err != nil {
			return false, err
		}
	}

	switch ast.BooleanOperator {
	case rules.BooleanOperator("AND"):
		return leftVal && rightVal, nil
	case rules.BooleanOperator("OR"):
		return leftVal || rightVal, nil
	default:
		return false, fmt.Errorf("invalid boolean operator: %s", ast.BooleanOperator)
	}
}

func (engine *RulesEngine) evaluateExpression(expression *rules.ConditionExpression, values map[string]string) (bool, error) {
	sensor, err := engine.database.GetSensor(expression.DeviceId, expression.SensorId)
	if err != nil {
		return false, err
	}

	switch sensor.DataType {
	case models.DataTypeBool:
		return engine.evaluateBoolExpression(expression, values)
	case models.DataTypeInt:
		return engine.evaluateIntExpression(expression, values)
	case models.DataTypeFloat:
		return engine.evaluateFloatExpression(expression, values)
	case models.DataTypeString:
		return engine.evaluateStringExpression(expression, values)
	default:
		return false, fmt.Errorf("unknown data type: %s", sensor.DataType)
	}
}

func (engine *RulesEngine) evaluateStringExpression(expression *rules.ConditionExpression, values map[string]string) (bool, error) {
	key := fmt.Sprintf("%s.%s.%s", expression.DeviceId, expression.SensorId, expression.Variable)
	value, ok := values[key]
	if !ok {
		return false, fmt.Errorf("unknown sensor value: %s", key)
	}

	switch expression.Operator {
	case rules.Operator("=="):
		return value == expression.Value, nil
	case rules.Operator("!="):
		return value != expression.Value, nil
	default:
		return false, fmt.Errorf("invalid operator for type string: %s", expression.Operator)
	}
}

func (engine *RulesEngine) evaluateBoolExpression(expression *rules.ConditionExpression, values map[string]string) (bool, error) {
	key := fmt.Sprintf("%s.%s.%s", expression.DeviceId, expression.SensorId, expression.Variable)
	value, ok := values[key]
	if !ok {
		return false, fmt.Errorf("unknown sensor value: %s", key)
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid value for type bool: %s", value)
	}

	expValue, err := strconv.ParseBool(expression.Value)
	if err != nil {
		return false, fmt.Errorf("invalid value for type bool: %s", expression.Value)
	}

	switch expression.Operator {
	case rules.Operator("=="):
		return boolValue == expValue, nil
	case rules.Operator("!="):
		return boolValue != expValue, nil
	default:
		return false, fmt.Errorf("invalid operator for type string: %s", expression.Operator)
	}
}

func (engine *RulesEngine) evaluateFloatExpression(expression *rules.ConditionExpression, values map[string]string) (bool, error) {
	key := fmt.Sprintf("%s.%s.%s", expression.DeviceId, expression.SensorId, expression.Variable)
	value, ok := values[key]
	if !ok {
		return false, fmt.Errorf("unknown sensor value: %s", key)
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false, fmt.Errorf("invalid value for type int: %s", value)
	}

	expValue, err := strconv.ParseFloat(expression.Value, 64)
	if err != nil {
		return false, fmt.Errorf("invalid value for type int: %s", expression.Value)
	}

	switch expression.Operator {
	case rules.Operator("=="):
		return floatVal == expValue, nil
	case rules.Operator("!="):
		return floatVal != expValue, nil
	case rules.Operator(">"):
		return floatVal > expValue, nil
	case rules.Operator("<"):
		return floatVal < expValue, nil
	case rules.Operator(">="):
		return floatVal >= expValue, nil
	case rules.Operator("<="):
		return floatVal <= expValue, nil
	default:
		return false, fmt.Errorf("invalid operator for type float: %s", expression.Operator)
	}
}

func (engine *RulesEngine) evaluateIntExpression(expression *rules.ConditionExpression, values map[string]string) (bool, error) {
	key := fmt.Sprintf("%s.%s.%s", expression.DeviceId, expression.SensorId, expression.Variable)
	value, ok := values[key]
	if !ok {
		return false, fmt.Errorf("unknown sensor value: %s", key)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return false, fmt.Errorf("invalid value for type int: %s", value)
	}

	expValue, err := strconv.Atoi(expression.Value)
	if err != nil {
		return false, fmt.Errorf("invalid value for type int: %s", expression.Value)
	}

	switch expression.Operator {
	case rules.Operator("=="):
		return intValue == expValue, nil
	case rules.Operator("!="):
		return intValue != expValue, nil
	case rules.Operator(">"):
		return intValue > expValue, nil
	case rules.Operator("<"):
		return intValue < expValue, nil
	case rules.Operator(">="):
		return intValue >= expValue, nil
	case rules.Operator("<="):
		return intValue <= expValue, nil
	default:
		return false, fmt.Errorf("invalid operator for type int: %s", expression.Operator)
	}
}

func (engine *RulesEngine) readDependentValues(deps []UsedSensorValue) (map[string]string, error) {
	results := make(map[string]string)

	for _, dep := range deps {
		key := fmt.Sprintf("%s.%s.%s", dep.DeviceId, dep.SensorId, dep.Type)
		if dep.Type == CurrentSensorValue {
			value, err := engine.database.GetCurrentSensorValue(dep.DeviceId, dep.SensorId)
			if err != nil {
				return nil, err
			}
			results[key] = value.Value
		} else if dep.Type == PreviousSensorValue {
			value, err := engine.database.GetPreviousSensorValue(dep.DeviceId, dep.SensorId)
			if err != nil {
				return nil, err
			}
			results[key] = value.Value
		} else {
			return nil, fmt.Errorf("unknown sensor value type: %s", dep.Type)
		}
	}

	return results, nil
}

func buildLookupTable(database RulesDatabase) (map[string][]rules.Rule, error) {
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
	ast, err := rule.ReadConditionAst()
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
