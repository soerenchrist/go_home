package background

import (
	"fmt"
	"log"
	"time"

	"github.com/soerenchrist/go_home/internal/db"
	"github.com/soerenchrist/go_home/internal/sensor"
	"github.com/soerenchrist/go_home/internal/value"
)

func PollSensorValues(database db.Database, outputBinding *value.OutputBindings) {
	sensors, err := database.ListPollingSensors()
	if err != nil {
		panic(err)
	}

	lastPolls := createPollingMap(sensors)
	channel := make(chan sensor.Sensor, 10)

	go readChannel(database, channel, outputBinding)

	log.Printf("Found %d sensors to poll\n", len(sensors))

	for {
		currentTimestamp := time.Now().Unix()
		for _, sensor := range sensors {
			lastPoll := lastPolls[sensor.ID]
			if currentTimestamp-lastPoll >= int64(sensor.PollingInterval) {
				channel <- sensor
				lastPolls[sensor.ID] = currentTimestamp
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func readChannel(database db.Database, channel chan sensor.Sensor, outputBinding *value.OutputBindings) {
	for {
		sensor := <-channel

		strategy, err := getStrategy(&sensor)
		if err != nil {
			log.Printf("Error polling sensor %s: %s\n", sensor.ID, err.Error())
			continue
		}
		result, err := strategy.PerformRequest(&sensor)
		if err != nil {
			log.Printf("Error polling sensor %s: %s\n", sensor.ID, err.Error())
			continue
		}
		err = database.AddSensorValue(result)
		if err != nil {
			log.Printf("Failed to save polling result %s: %s\n", sensor.ID, err.Error())
			continue
		}
		outputBinding.Push(*result)
		log.Printf("Polled sensor %s successfully with result %s\n", sensor.ID, result.Value)
	}
}

func getStrategy(s *sensor.Sensor) (RequestStrategy, error) {
	switch s.PollingStrategy {
	case sensor.PollingStrategyPing:
		return &PingStrategy{}, nil
	}

	return nil, fmt.Errorf("unknown polling strategy %s", s.PollingStrategy)
}

func createPollingMap(sensors []sensor.Sensor) map[string]int64 {
	pollingMap := make(map[string]int64)

	for _, sensor := range sensors {
		pollingMap[sensor.ID] = 0
	}

	return pollingMap
}
