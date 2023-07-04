package output_test

import (
	"testing"

	"github.com/soerenchrist/go_home/pkg/output"
)

func TestOutputBindingShouldReceiveValue(t *testing.T) {
	manager := output.NewManager()

	outputBinding := output.NewChannelOutput()
	manager.Register(outputBinding)

	manager.Push(output.BindingValue{Value: "10"})

	rec := <-outputBinding.Channel

	if rec.Value != "10" {
		t.Fatalf("Failed to receive value")
	}
}

func TestOutputBindingWithMultipleOutputs(t *testing.T) {
	manager := output.NewManager()
	bindings := make([]*output.ChannelOutputBinding, 0)

	for i := 0; i < 10; i++ {
		outputBinding := output.NewChannelOutput()
		manager.Register(outputBinding)
		bindings = append(bindings, outputBinding)
	}

	manager.Push(output.BindingValue{Value: "10"})

	for _, b := range bindings {
		select {
		case <-b.Channel:
			t.Log("Received value")
		default:
			t.Fatalf("Failed to receive value")
		}
	}
}
