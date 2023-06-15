package templating_test

import (
	"io"
	"testing"

	"github.com/soerenchrist/mini_home/templating"
)

var template = `{
		"device": "{{.device_id}}",
		"command": "{{.command_id}}",
		"payload": "{{.p_payload}},
		"now": "{{.now}}"{{if .p_show_device_name}},
		"device_name": "{{.p_device_name}}"{{- end}}
		}`

func TestPrepareTemplate(t *testing.T) {

	var params templating.TemplateParameters = make(map[string]string)
	params["device_id"] = "device_id"
	params["command_id"] = "command_id"
	params["p_payload"] = "payload"
	params["now"] = "now"

	reader, err := templating.PrepareCommandTemplate(template, &params)
	if err != nil {
		t.Errorf("Error preparing template: %s", err)
	}

	expected := `{
		"device": "device_id",
		"command": "command_id",
		"payload": "payload,
		"now": "now"
		}`

	if reader == nil {
		t.Errorf("Reader is nil")
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("Error reading from reader: %s", err)
	}

	if string(bytes) != expected {
		t.Errorf("Expected %s, got %s", expected, string(bytes))
	}
}

func TestPrepareTemplateWithIf(t *testing.T) {

	var params templating.TemplateParameters = make(map[string]string)
	params["device_id"] = "device_id"
	params["command_id"] = "command_id"
	params["p_payload"] = "payload"
	params["now"] = "now"
	params["p_device_name"] = "device_name"
	params["p_show_device_name"] = "true"

	reader, err := templating.PrepareCommandTemplate(template, &params)
	if err != nil {
		t.Errorf("Error preparing template: %s", err)
	}

	expected := `{
		"device": "device_id",
		"command": "command_id",
		"payload": "payload,
		"now": "now",
		"device_name": "device_name"
		}`

	if reader == nil {
		t.Errorf("Reader is nil")
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("Error reading from reader: %s", err)
	}

	if string(bytes) != expected {
		t.Errorf("Expected %s, got %s", expected, string(bytes))
	}
}
