package background

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-ping/ping"
	"github.com/soerenchrist/go_home/internal/sensor"
	"github.com/soerenchrist/go_home/internal/value"
)

type RequestStrategy interface {
	PerformRequest(sensor *sensor.Sensor) (*value.SensorValue, error)
}

type PingStrategy struct{}

func (strgy *PingStrategy) PerformRequest(s *sensor.Sensor) (*value.SensorValue, error) {
	if s.DataType != sensor.DataTypeBool {
		return nil, fmt.Errorf("PingStrategy can only be used with boolean data types")
	}

	pinger, err := ping.NewPinger(s.PollingEndpoint)
	if err != nil {
		return nil, err
	}

	pinger.Count = 3
	pinger.Timeout = 5 * time.Second
	err = pinger.Run()
	if err != nil {
		return nil, err
	}

	stats := pinger.Statistics()
	log.Debugf("Ping %s: %v\n", s.PollingEndpoint, stats)

	reachable := stats.PacketsRecv > 0

	return &value.SensorValue{
		Value:     strconv.FormatBool(reachable),
		SensorID:  s.ID,
		DeviceID:  s.DeviceID,
		Timestamp: time.Now(),
	}, nil
}
