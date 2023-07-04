package background

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
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

	log.Debug().Int("count", len(sensors)).Msgf("Found %d sensors to poll\n", len(sensors))

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
			log.Error().Err(err).Str("sensor_id", sensor.ID).Msg("Error polling sensor")
			continue
		}
		result, err := strategy.PerformRequest(&sensor)
		if err != nil {
			log.Error().Err(err).Str("sensor_id", sensor.ID).Msg("Error polling sensor")
			continue
		}
		err = database.AddSensorValue(result)
		if err != nil {
			log.Error().Err(err).Str("sensor_id", sensor.ID).Msg("Failed to save polling result")
			continue
		}
		outputBinding.Push(*result)
		log.Debug().
			Str("sensor_id", sensor.ID).
			Str("polling_result", result.Value).
			Msg("Polled sensor successfully")
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
