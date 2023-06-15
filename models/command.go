package models

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

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
	body, err := c.prepareBody(device, params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(c.Method, c.Endpoint, body)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func (c *Command) prepareBody(device *Device, params *CommandParameters) (io.Reader, error) {
	if len(c.PayloadTemplate) == 0 {
		return nil, nil
	}

	t := template.Must(template.New("payload").Parse(c.PayloadTemplate))

	var data CommandParameters = make(map[string]string)
	data["command_id"] = c.ID
	data["command_name"] = c.Name
	data["device_id"] = device.ID
	data["device_name"] = device.Name
	data["now"] = util.GetTimestamp()

	for key, value := range *params {
		data[fmt.Sprintf("p_%s", key)] = value
	}

	var b bytes.Buffer

	err := t.Execute(&b, data)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(b.String()), nil
}

type InvocationResult struct {
	Response   string `json:"response"`
	StatusCode int    `json:"statusCode"`
}
