package value

import (
	"log"
)

type OutputBindings struct {
	channels []chan SensorValue
}

func NewOutputBindings() *OutputBindings {
	channels := make([]chan SensorValue, 0)

	return &OutputBindings{
		channels: channels,
	}
}

func (bindings *OutputBindings) Register(channel chan SensorValue) {
	bindings.channels = append(bindings.channels, channel)
}

func (bindings *OutputBindings) Push(val SensorValue) {
	for _, channel := range bindings.channels {
		bindings.send(channel, val)
	}
}

func (bindings *OutputBindings) send(channel chan SensorValue, val SensorValue) {
	select {
	case channel <- val:
		log.Println("Sent value to output binding")
	default:
		log.Println("Could not sent to output binding")
	}
}
