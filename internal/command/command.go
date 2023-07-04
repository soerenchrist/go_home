package command

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/soerenchrist/go_home/internal/device"
	"github.com/soerenchrist/go_home/internal/util"
)

type CommandParameters map[string]string

type Command struct {
	ID              string `json:"id"`
	DeviceID        string `json:"device_id"`
	Name            string `json:"name"`
	PayloadTemplate string `json:"payload"`
	Endpoint        string `json:"endpoint"`
	Method          string `json:"method"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Command) String() string {
	return fmt.Sprintf("Command<%s %s>", c.ID, c.Name)
}

func (c *Command) Invoke(device *device.Device, params *CommandParameters) (*http.Response, error) {
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

func (command *Command) prepareBody(device *device.Device, params *CommandParameters) (io.Reader, error) {
	if len(command.PayloadTemplate) == 0 {
		return nil, nil
	}

	var data TemplateParameters = make(map[string]string)
	data["command_id"] = command.ID
	data["command_name"] = command.Name
	data["device_id"] = device.ID
	data["device_name"] = device.Name
	data["now"] = util.GetTimestamp()

	for key, value := range *params {
		data[fmt.Sprintf("p_%s", key)] = value
	}

	return PrepareCommandTemplate(command.PayloadTemplate, &data)
}

type InvocationResult struct {
	Response   string `json:"response"`
	StatusCode int    `json:"statusCode"`
}

type TemplateParameters map[string]string

func PrepareCommandTemplate(payloadTemplate string, params *TemplateParameters) (io.Reader, error) {

	t := template.Must(template.New("payload").Parse(payloadTemplate))
	var b bytes.Buffer

	err := t.Execute(&b, params)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(b.String()), nil
}

type CreateCommandRequest struct {
	Name            string `json:"name"`
	PayloadTemplate string `json:"payload_template"`
	Endpoint        string `json:"endpoint"`
	Method          string `json:"method"`
}
