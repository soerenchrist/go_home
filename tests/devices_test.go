package tests

import (
	"encoding/json"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
)

func TestGetDevices(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices")

	assert.Equal(t, 200, w.Code)

	var data []models.Device
	err := json.Unmarshal(w.Body.Bytes(), &data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(data))
}

func TestCreateDevice_ShouldReturnTheCreatedDevice_WhenTheBodyIsValid(t *testing.T) {
	body := `{
		"name": "Test Device"
	}`
	w := RecordPostCall(t, "/api/v1/devices", body)

	assert.Equal(t, w.Code, 201)

	var data models.Device
	err := json.Unmarshal(w.Body.Bytes(), &data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Test Device", data.Name)
	assert.Equal(t, IsValidUuid(data.ID), true)
}

func TestCreateDevice_ShouldReturn400_WhenBodyIsInvalid(t *testing.T) {
	body := `{
		"name": "Test Device"
	`
	w := RecordPostCall(t, "/api/v1/devices", body)

	assert.Equal(t, w.Code, 400)
}

func TestCreateDevice_ShouldReturn400_WhenNameIsNotValid(t *testing.T) {
	body := `{
		"name": "XX"
	}
	`
	w := RecordPostCall(t, "/api/v1/devices", body)

	assert.Equal(t, w.Code, 400)
}

func TestCreateDevice_ShouldBeInDatabase_WhenDataIsValid(t *testing.T) {
	body := `{
		"name": "TestDevice"
	}
	`

	// check, if the device is in the database
	validateDb := func(database db.Database) {
		devices, err := database.ListDevices()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(devices))
		assert.Equal(t, "TestDevice", devices[2].Name)
		assert.Equal(t, IsValidUuid(devices[2].ID), true)
	}

	w := RecordPostCallWithDb(t, "/api/v1/devices", body, validateDb)

	assert.Equal(t, w.Code, 201)
}

func TestGetDevice_ShouldReturn404_WhenTheGivenIdDoesNotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/123")

	assert.Equal(t, w.Code, 404)
}

func TestGetDevice_ShouldReturnDevice_WhenTheGivenIdDoesExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1")

	assert.Equal(t, w.Code, 200)
	var data models.Device
	err := json.Unmarshal(w.Body.Bytes(), &data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1", data.ID)
	assert.Equal(t, "My Device 1", data.Name)
}

func TestDeleteDevice_ShouldReturn404_WhenTheGivenIdDoesNotExist(t *testing.T) {
	w := RecordDeleteCall(t, "/api/v1/devices/123")

	assert.Equal(t, w.Code, 404)
}

func TestDeleteDevice_ShouldReturn204AndDelete_WhenTheGivenIdDoesExist(t *testing.T) {
	validator := func(database db.Database) {
		devices, err := database.ListDevices()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 1, len(devices))
		assert.Equal(t, "My Device 2", devices[0].Name)
		assert.Equal(t, "2", devices[0].ID)
	}

	w := RecordDeleteCallWithDb(t, "/api/v1/devices/1", validator)

	assert.Equal(t, w.Code, 204)
}
