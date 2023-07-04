package output

import (
	"time"

	"github.com/rs/zerolog/log"
)

type BindingValue struct {
	Timestamp time.Time
	Value     string
	SensorID  string
	DeviceID  string
}

type OutputBinding interface {
	Handle(value BindingValue)
}

type OutputBindingsManager struct {
	bindings []OutputBinding
}

type ChannelOutputBinding struct {
	Channel chan BindingValue
}

func NewChannelOutput() *ChannelOutputBinding {
	return &ChannelOutputBinding{
		Channel: make(chan BindingValue, 2),
	}
}

func (b ChannelOutputBinding) Handle(value BindingValue) {
	select {
	case b.Channel <- value:
	default:
	}
}

func NewManager() *OutputBindingsManager {
	bindings := make([]OutputBinding, 0)

	return &OutputBindingsManager{
		bindings: bindings,
	}
}

func (m *OutputBindingsManager) Register(binding OutputBinding) {
	m.bindings = append(m.bindings, binding)
}

func (m *OutputBindingsManager) Push(val BindingValue) {
	for _, channel := range m.bindings {
		m.send(channel, val)
	}
}

func (m *OutputBindingsManager) send(binding OutputBinding, val BindingValue) {
	log.Debug().
		Str("value", val.Value).
		Str("sensor_id", val.SensorID).
		Str("device_id", val.DeviceID).
		Msg("Sending value to output binding")
	binding.Handle(val)
}
