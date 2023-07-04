package evaluation

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/soerenchrist/go_home/internal/command"
	"github.com/soerenchrist/go_home/internal/rules"
	"github.com/soerenchrist/go_home/internal/sensor"
	"github.com/soerenchrist/go_home/internal/value"
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
	database    rules.RulesDatabase
	lookupTable map[string][]rules.Rule
}

func NewRulesEngine(database rules.RulesDatabase) *RulesEngine {
	lookupTable, err := buildLookupTable(database)
	if err != nil {
		panic(err)
	}
	return &RulesEngine{lookupTable: lookupTable, database: database}
}

func (engine *RulesEngine) ListenForValues(sensorsChannel chan value.SensorValue) {
	log.Debug().Msg("Listening for sensor values...")
	for {
		sensor := <-sensorsChannel

		key := sensor.DeviceID + "." + sensor.SensorID
		if rules, ok := engine.lookupTable[key]; ok {
			for _, rule := range rules {
				log.Debug().Int64("rule_id", rule.Id).Str("rule_name", rule.Name).Msg("Evaluating rule")
				evalResult, err := engine.EvaluateRule(&rule)
				if err != nil {
					log.Error().Err(err).Msg("Error evaluating rule")
					continue
				}
				log.Debug().Str("rule_name", rule.Name).Bool("eval_result", evalResult).Msgf("Rule '%s' evaluated to %t", rule.Name, evalResult)
				if evalResult {
					err := engine.executeRule(rule)
					if err != nil {
						log.Error().Err(err).Msg("Error executing rule")
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

	cmd, err := engine.database.GetCommand(action.DeviceId, action.CommandId)
	if err != nil {
		return fmt.Errorf("error reading command: %v", err)
	}

	log.Debug().Str("command_id", cmd.ID).Str("device_id", cmd.DeviceID).Msg("Executing command")

	var params command.CommandParameters
	if action.Payload != "" {
		err := json.Unmarshal([]byte(action.Payload), &params)
		if err != nil {
			return err
		}
	}

	resp, err := cmd.Invoke(device, &params)
	if err != nil {
		return fmt.Errorf("error invoking command: %v", err)
	}

	log.Debug().Int("response_status", resp.StatusCode).Msgf("Command response status: %d \n", resp.StatusCode)
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
		log.Debug().Str("sensor_name", key).Str("sensor_value", value).Msgf("Sensor value: %s = %v\n", key, value)
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
	s, err := engine.database.GetSensor(expression.DeviceId, expression.SensorId)
	if err != nil {
		return false, err
	}

	switch s.DataType {
	case sensor.DataTypeBool:
		return engine.evaluateBoolExpression(expression, values)
	case sensor.DataTypeInt:
		return engine.evaluateIntExpression(expression, values)
	case sensor.DataTypeFloat:
		return engine.evaluateFloatExpression(expression, values)
	case sensor.DataTypeString:
		return engine.evaluateStringExpression(expression, values)
	default:
		return false, fmt.Errorf("unknown data type: %s", s.DataType)
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

func buildLookupTable(database rules.RulesDatabase) (map[string][]rules.Rule, error) {
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
