package tests

import (
	"encoding/json"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/soerenchrist/go_home/internal/db"
	"github.com/soerenchrist/go_home/internal/models"
)

func TestListCommands_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/123/commands")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestListCommands_ShouldReturnCommands(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1/commands")

	assert.Equal(t, w.Code, 200)

	var commands []models.Command
	err := json.Unmarshal(w.Body.Bytes(), &commands)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(commands), 1)
	assert.Equal(t, commands[0].ID, "C1")
	assert.Equal(t, commands[0].Name, "Turn on")
}

func TestGetCommand_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/123/commands/C1")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestGetCommand_ShouldReturn404_WhenCommandDoesNotExist(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1/commands/C2")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Command not found")
}

func TestGetCommand_ShouldReturnCommand(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/devices/1/commands/C1")

	assert.Equal(t, w.Code, 200)

	var command models.Command
	err := json.Unmarshal(w.Body.Bytes(), &command)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, command.ID, "C1")
	assert.Equal(t, command.Name, "Turn on")
}

func TestCreateCommand_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	body := `{
		"name": "Turn on"
	}`
	w := RecordPostCall(t, "/api/v1/devices/123/commands", body)

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestCreateCommand_ShouldReturn400_WhenBodyIsMalformed(t *testing.T) {
	body := `{
		"name": "Turn on"
	`
	w := RecordPostCall(t, "/api/v1/devices/1/commands", body)

	assert.Equal(t, w.Code, 400)
}

func TestCreateCommand_ShouldReturn400_WhenBodyIsInvalid(t *testing.T) {
	bodies := []string{
		`{}`,
		`{"name": ""}`,
		`{"name": "Test"}`,
		`{"name": "Test", "endpoint": "http://localhost:8080"}`,
		`{"name": "Test", "endpoint": "http://localhost:8080", "payload_template": "on"}`,
		`{"name": "Test", "endpoint": "http://localhost:8080", "payload_template": "on", "method": "TEST"}`,
	}
	messages := []string{
		"Name is required",
		"Name is required",
		"Endpoint is required",
		"Payload template is required",
		"Method must be one of GET, POST, PUT or DELETE",
		"Method must be one of GET, POST, PUT or DELETE",
	}

	for i, body := range bodies {
		w := RecordPostCall(t, "/api/v1/devices/1/commands", body)
		expectedMessage := messages[i]

		assert.Equal(t, w.Code, 400)
		t.Log(body)
		assertErrorMessageEquals(t, w.Body.Bytes(), expectedMessage)
	}
}

func TestCreateCommand_ShouldAddCommandToDatabase(t *testing.T) {
	body := `{
		"name": "Test",
		"endpoint": "http://localhost:8080",
		"payload_template": "on",
		"method": "POST"
	}`

	validator := func(database db.Database) {
		commands, err := database.ListCommands("1")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(commands), 2)
		assert.Equal(t, commands[1].Name, "Test")
		assert.Equal(t, commands[1].Endpoint, "http://localhost:8080")
		assert.Equal(t, commands[1].PayloadTemplate, "on")
		assert.Equal(t, commands[1].Method, "POST")
	}

	w := RecordPostCallWithDb(t, "/api/v1/devices/1/commands", body, validator)

	assert.Equal(t, w.Code, 201)

	var command models.Command
	err := json.Unmarshal(w.Body.Bytes(), &command)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, command.Name, "Test")
	assert.Equal(t, command.Endpoint, "http://localhost:8080")
	assert.Equal(t, command.PayloadTemplate, "on", "Response should contain payload template")
	assert.Equal(t, command.Method, "POST")
}

func TestDeleteCommand_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	w := RecordDeleteCall(t, "/api/v1/devices/123/commands/C1")

	assert.Equal(t, w.Code, 404)
	t.Log(w.Body.String())
	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestDeleteCommand_ShouldReturn404_WhenCommandDoesNotExist(t *testing.T) {
	w := RecordDeleteCall(t, "/api/v1/devices/1/commands/C2")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Command not found")
}

func TestDeleteCommand_ShouldDeleteCommand(t *testing.T) {
	validator := func(database db.Database) {
		commands, err := database.ListCommands("1")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(commands), 0)
	}

	w := RecordDeleteCallWithDb(t, "/api/v1/devices/1/commands/C1", validator)

	assert.Equal(t, w.Code, 204)
}

func TestInvokeCommand_ShouldReturn404_WhenDeviceDoesNotExist(t *testing.T) {
	w := RecordPostCall(t, "/api/v1/devices/123/commands/C1/invoke", "")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Device not found")
}

func TestInvokeCommand_ShouldReturn404_WhenCommandDoesNotExist(t *testing.T) {
	w := RecordPostCall(t, "/api/v1/devices/1/commands/C2/invoke", "")

	assert.Equal(t, w.Code, 404)
	assertErrorMessageEquals(t, w.Body.Bytes(), "Command not found")
}
