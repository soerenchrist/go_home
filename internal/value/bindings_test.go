package value

import (
	"testing"
)

func TestOutputBindingShouldReceiveValue(t *testing.T) {
	output := NewOutputBindings()

	binding := make(chan SensorValue, 10)

	output.Register(binding)

	output.Push(SensorValue{Value: "10"})

	rec := <-binding

	if rec.Value != "10" {
		t.Fatalf("Failed to receive value")
	}
}

func TestOutputBindingWithMultipleOutputs(t *testing.T) {
	output := NewOutputBindings()
	channels := make([]chan SensorValue, 0)

	for i := 0; i < 10; i++ {
		c := make(chan SensorValue, 10)
		channels = append(channels, c)
		output.Register(c)

	}

	output.Push(SensorValue{Value: "10"})

	for _, c := range channels {
		select {
		case <-c:
		default:
			t.Fatalf("Failed to receive value")
		}
	}
}
