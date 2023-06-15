package tests

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/soerenchrist/mini_home/db"
	"github.com/soerenchrist/mini_home/models"
	"github.com/soerenchrist/mini_home/server"
	"github.com/soerenchrist/mini_home/util"
)

func TestAddSensorValue_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	body := `{
		"value": 1.23
	}`
	w := RecordPostCall(t, "/api/v1/devices/123/sensors/1/values", body)

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestAddSensorValue_ShouldReturn404_WhenSensorDoesNotExist(t *testing.T) {
	body := `{
		"value": 1.23
	}`
	w := RecordPostCall(t, "/api/v1/devices/1/sensors/123/values", body)

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Sensor not found")
}

func TestAddSensorValue_ShouldReturn400_WhenBodyIsMalformed(t *testing.T) {
	body := `{
		"value": "1.23
	}`
	w := RecordPostCall(t, "/api/v1/devices/1/sensors/S1/values", body)

	assert.Equal(t, w.Code, 400)
}

func TestAddSensorValue_ShouldReturn400_WhenValueIsWrongDataType(t *testing.T) {
	bodies := []string{
		`{"value": "NOT_A_FLOAT"}`,
		`{"value": "10"}`,
		`{"value": "10.123"}`,
	}

	sensors := []struct {
		device string
		sensor string
	}{
		{device: "1", sensor: "S1"},
		{device: "1", sensor: "S2"},
		{device: "2", sensor: "S3"},
	}

	messages := []string{
		"Sensor value is not a float",
		"Sensor value is not a bool",
		"Sensor value is not an int",
	}

	for i, body := range bodies {
		sensor := sensors[i]

		url := fmt.Sprintf("/api/v1/devices/%s/sensors/%s/values", sensor.device, sensor.sensor)

		w := RecordPostCall(t, url, body)

		assert.Equal(t, w.Code, 400)

		expectedMessage := messages[i]
		assertErrorMessageEquals(t, w.Body.Bytes(), expectedMessage)
	}
}

func TestAddSensorValue_ShouldReturn400_WhenSensorIsPolling(t *testing.T) {
	body := `{
		"value": "true"
	}`
	w := RecordPostCall(t, "/api/v1/devices/1/sensors/S2/values", body)

	assert.Equal(t, w.Code, 400)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Sending values to a polling sensor is not allowed")
}

func TestAddSensorValue_ShouldReturn400_WhenTimestampIsInvalid(t *testing.T) {
	body := `{
		"value": "1.23",
		"timestamp": "Something" 
	}`

	w := RecordPostCall(t, "/api/v1/devices/1/sensors/S1/values", body)

	assert.Equal(t, w.Code, 400)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Something is not a valid RFC3339 timestamp")
}

func TestAddSensorValue_ShouldAddSensorValueToDb(t *testing.T) {
	body := `{
		"value": "1.23"
	}`
	validator := func(database db.DevicesDatabase) {
		value, err := database.GetCurrentSensorValue("1", "S1")
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, value.Value, "1.23")
		assert.Equal(t, util.ValidateTimestamp(value.Timestamp), nil)
	}

	w := RecordPostCallWithDb(t, "/api/v1/devices/1/sensors/S1/values", body, validator)

	assert.Equal(t, w.Code, 201)

	var result models.SensorValue
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, result.Value, "1.23")
	assert.Equal(t, util.ValidateTimestamp(result.Timestamp), nil)
}

func TestGetCurrentSensorValue_ShouldReturn404_WhenDeviceDoesnotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/123/sensors/1/current")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestGetCurrentSensorValue_ShouldReturn404_WhenSensorDoesnotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1/sensors/123/current")

	assert.Equal(t, w.Code, 404)
	t.Log(w.Body.String())
	assertErrorMessageEquals(t, w.Body.Bytes(), "Sensor not found")
}

func TestGetcurrentSensorValue_ShouldReturn404_WhenNoValueExists(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1/sensors/S1/current")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "No sensor value found")
}

func TestGetCurrentSensorValue_ShouldReturnSensorValue_WhenOneExists(t *testing.T) {
	w := httptest.NewRecorder()
	filename := t.Name() + ".db"
	database := CreateTestDatabase(filename)
	defer CloseTestDatabase(database, filename)
	timestamp := "2019-01-01T00:00:00Z"
	err := database.AddSensorValue(&models.SensorValue{SensorID: "S1", Value: "1.23", DeviceID: "1", Timestamp: timestamp})
	if err != nil {
		t.Error(err)
	}
	router := server.NewRouter(database)

	req := httptest.NewRequest("GET", "/api/v1/devices/1/sensors/S1/current", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)

	var result models.SensorValue
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, result.Value, "1.23")
	assert.Equal(t, result.Timestamp, timestamp)
}

func TestGetSensorValues_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/123/sensors/1/values")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestGetSensorValues_ShouldReturn404_WhenSensorDoesNotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1/sensors/123/values")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Sensor not found")
}

func TestGetSensorValues_ShouldReturnEmptyList_WhenLastValuesIsTooLongAgo(t *testing.T) {
	w := httptest.NewRecorder()
	filename := t.Name() + ".db"
	database := CreateTestDatabase(filename)
	defer CloseTestDatabase(database, filename)

	timestampOneHourAgo := getTimestamp(time.Now().Add(-1*time.Hour - 1*time.Minute))

	err := database.AddSensorValue(&models.SensorValue{SensorID: "S1", Value: "1.23", DeviceID: "1", Timestamp: timestampOneHourAgo})
	if err != nil {
		t.Error(err)
	}
	router := server.NewRouter(database)

	req := httptest.NewRequest("GET", "/api/v1/devices/1/sensors/S1/values", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "[]")
}

func TestGetSensorValues_ShouldReturnValues_WhenValuesAreInTimeFrame(t *testing.T) {

	w := httptest.NewRecorder()
	filename := t.Name() + ".db"
	database := CreateTestDatabase(filename)
	defer CloseTestDatabase(database, filename)

	timestampOneHourAgo := getTimestamp(time.Now().Add(-1*time.Hour - 1*time.Minute))

	err := database.AddSensorValue(&models.SensorValue{SensorID: "S1", Value: "1.23", DeviceID: "1", Timestamp: timestampOneHourAgo})
	if err != nil {
		t.Error(err)
	}
	router := server.NewRouter(database)
	req := httptest.NewRequest("GET", "/api/v1/devices/1/sensors/S1/values?timeframe=2h", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)

	var result []models.SensorValue
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Value, "1.23")
	assert.Equal(t, result[0].Timestamp, timestampOneHourAgo)
}

func getTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}
