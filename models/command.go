package models

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/soerenchrist/mini_home/util"
)

type CommandParameters map[string]string

type Command struct {
	ID              string `json:"id"`
	DeviceID        string `json:"device_id"`
	Name            string `json:"name"`
	PayloadTemplate string `json:"payload"`
	Endpoint        string `json:"endpoint"`
	Method          string `json:"method"`
}

func (c *Command) String() string {
	return fmt.Sprintf("Command<%s %s>", c.ID, c.Name)
}

func (c *Command) Invoke(device *Device, params *CommandParameters) (*http.Response, error) {
	req, err := http.NewRequest(c.Method, c.Endpoint, c.prepareBody(device, params))
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func (c *Command) prepareBody(device *Device, params *CommandParameters) io.Reader {
	if len(c.PayloadTemplate) == 0 {
		return nil
	}

	template := strings.ReplaceAll(c.PayloadTemplate, "${command_id}", c.ID)
	template = strings.ReplaceAll(template, "${command_name}", c.Name)
	template = strings.ReplaceAll(template, "${device_id}", device.ID)
	template = strings.ReplaceAll(template, "${device_name}", device.Name)
	template = strings.ReplaceAll(template, "${now}", util.GetTimestamp())

	for key, value := range *params {
		template = strings.ReplaceAll(template, fmt.Sprintf("${p_%s}", key), value)
	}

	return strings.NewReader(template)
}
