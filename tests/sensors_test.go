package tests

import (
	"encoding/json"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
)

func TestGetSensors_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/123/sensors")

	assert.Equal(t, w.Code, 404)

	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestGetSensors_ShouldReturnSensors_WhenDeviceDoesExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1/sensors")

	assert.Equal(t, w.Code, 200)

	var data []models.Sensor
	err := json.Unmarshal(w.Body.Bytes(), &data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(data))
}

func TestCreateSensor_ShouldReturn400_WhenBodyIsInvalid(t *testing.T) {
	body := `{
		"id": "test_sensor",
		"name": "Test Sensor"
	`
	w := RecordPostCall(t, "/api/v1/devices/1/sensors", body)

	assert.Equal(t, w.Code, 400)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Invalid JSON")
}

func TestCreateSensor_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	body := `{
		"id": "test_sensor",
		"name": "Test Sensor"
	}
	`
	w := RecordPostCall(t, "/api/v1/devices/123/sensors", body)

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestCreateSensor_ShouldReturn400_WhenSensorDoesAlreadyExist(t *testing.T) {
	body := `{
		"id": "S1",
		"name": "Test Sensor"
	}
	`
	w := RecordPostCall(t, "/api/v1/devices/1/sensors", body)

	assert.Equal(t, w.Code, 400)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Sensor with id S1 does already exist")
}

func TestCreateSensor_ShouldReturn400_WhenSensorIsInvalid(t *testing.T) {
	bodies := []string{
		`{
			"id": "my_sensor",
			"name": "T"
		}`,
		`{
			"id": "my_sensor",
			"name": "Test Sensor",
			"unit": "XX",
			"data_type": "bool"
		}`,
		`{
			"id": "my_sensor",
			"name": "Test Sensor",
			"unit": "XX",
			"data_type": "string"
		}`,
		`{
			"id": "my_sensor",
			"name": "Test Sensor",
			"data_type": "int",
			"type": "polling",
			"polling_interval": 0,
			"polling_strategy": "ping",
			"polling_endpoint": "http://"
		}`,
		`{
			"id": "my_sensor",
			"name": "Test Sensor",
			"data_type": "float",
			"type": "polling",
			"polling_interval": 10,
			"polling_strategy": "ping",
			"polling_endpoint": ""
		}`,
		`{
			"id": "my_sensor",
			"name": "Test Sensor",
			"data_type": "float",
			"type": "polling",
			"polling_interval": 10,
			"polling_strategy": "something",
			"polling_endpoint": "some_endpoint"
		}`,
	}
	expectedMessages := []string{
		"Name must be at least 3 characters long",
		"Unit is not allowed for this data type",
		"Unit is not allowed for this data type",
		"Polling interval must be greater than 0",
		"Polling endpoint is required",
		"Invalid polling strategy",
	}

	for i, body := range bodies {
		w := RecordPostCall(t, "/api/v1/devices/1/sensors", body)

		assert.Equal(t, w.Code, 400)
		expectedMessage := expectedMessages[i]

		assertErrorMessageEquals(t, w.Body.Bytes(), expectedMessage)
	}
}

func TestCreateSensor_ShouldAddSensorToDb_WhenBodyIsValid(t *testing.T) {
	body := `{
		"id": "my_sensor",
		"name": "Test Sensor",
		"data_type": "float",
		"type": "polling",
		"polling_interval": 10,
		"polling_strategy": "ping",
		"polling_endpoint": "http://"
	}`

	validator := func(database db.Database) {
		sensors, err := database.ListSensors("1")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(sensors))
		assert.Equal(t, "my_sensor", sensors[2].ID)
		assert.Equal(t, "Test Sensor", sensors[2].Name)
		assert.Equal(t, models.DataTypeFloat, sensors[2].DataType)
		assert.Equal(t, models.SensorTypePolling, sensors[2].Type)
		assert.Equal(t, 10, sensors[2].PollingInterval)
		assert.Equal(t, models.PollingStrategyPing, sensors[2].PollingStrategy)
		assert.Equal(t, "http://", sensors[2].PollingEndpoint)
	}

	w := RecordPostCallWithDb(t, "/api/v1/devices/1/sensors", body, validator)

	assert.Equal(t, w.Code, 201)

	var sensor models.Sensor
	err := json.Unmarshal(w.Body.Bytes(), &sensor)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "my_sensor", sensor.ID)
	assert.Equal(t, "Test Sensor", sensor.Name)
	assert.Equal(t, models.DataTypeFloat, sensor.DataType)
	assert.Equal(t, models.SensorTypePolling, sensor.Type)
	assert.Equal(t, 10, sensor.PollingInterval)
	assert.Equal(t, models.PollingStrategyPing, sensor.PollingStrategy)
	assert.Equal(t, "http://", sensor.PollingEndpoint)
}

func TestDeleteSensor_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	w := RecordDeleteCall(t, "/api/v1/devices/123/sensors/1")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Sensor not found")
}
func TestDeleteSensor_ShouldReturn404_WhenSensorDoesNotExist(t *testing.T) {
	w := RecordDeleteCall(t, "/api/v1/devices/1/sensors/S4")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Sensor not found")
}

func TestDeleteSensor_ShouldDeleteSensorFromDb_WhenSensorDoesExist(t *testing.T) {
	validator := func(database db.Database) {
		sensors, err := database.ListSensors("1")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 1, len(sensors))
	}

	w := RecordDeleteCallWithDb(t, "/api/v1/devices/1/sensors/S1", validator)

	assert.Equal(t, w.Code, 204)
}

func TestGetSensor_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/123/sensors/1")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Sensor not found")
}

func TestGetSensor_ShouldReturn404_WhenSensorDoesNotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1/sensors/1123")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Sensor not found")
}

func TestGetSensor_ShouldReturnSensor_WhenSensorDoesExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1/sensors/S1")

	assert.Equal(t, w.Code, 200)

	var sensor models.Sensor
	err := json.Unmarshal(w.Body.Bytes(), &sensor)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "S1", sensor.ID)
	assert.Equal(t, "Temperature", sensor.Name)
	assert.Equal(t, models.DataTypeFloat, sensor.DataType)
	assert.Equal(t, models.SensorTypeExternal, sensor.Type)
	assert.Equal(t, 0, sensor.PollingInterval)
}
