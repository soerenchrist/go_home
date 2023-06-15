package background

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-ping/ping"
	"github.com/soerenchrist/go_home/models"
	"github.com/soerenchrist/go_home/util"
)

type RequestStrategy interface {
	PerformRequest(sensor *models.Sensor) (*models.SensorValue, error)
}

type PingStrategy struct{}

func (s *PingStrategy) PerformRequest(sensor *models.Sensor) (*models.SensorValue, error) {
	if sensor.DataType != models.DataTypeBool {
		return nil, fmt.Errorf("PingStrategy can only be used with boolean data types")
	}

	pinger, err := ping.NewPinger(sensor.PollingEndpoint)
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
	log.Printf("Ping %s: %v\n", sensor.PollingEndpoint, stats)

	reachable := stats.PacketsRecv > 0

	return &models.SensorValue{
		Value:     strconv.FormatBool(reachable),
		SensorID:  sensor.ID,
		DeviceID:  sensor.DeviceID,
		Timestamp: util.GetTimestamp(),
	}, nil
}
