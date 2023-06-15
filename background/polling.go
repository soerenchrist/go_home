package background

import (
	"fmt"
	"log"
	"time"

	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
)

func PollSensorValues(database db.DevicesDatabase) {
	sensors, err := database.ListPollingSensors()
	if err != nil {
		panic(err)
	}

	lastPolls := createPollingMap(sensors)
	channel := make(chan models.Sensor)

	go readChannel(database, channel)

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

func readChannel(database db.DevicesDatabase, channel chan models.Sensor) {
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
		log.Printf("Polled sensor %s successfully with result %s\n", sensor.ID, result.Value)
	}
}

func getStrategy(sensor *models.Sensor) (RequestStrategy, error) {
	switch sensor.PollingStrategy {
	case models.PollingStrategyPing:
		return &PingStrategy{}, nil
	}

	return nil, fmt.Errorf("unknown polling strategy %s", sensor.PollingStrategy)
}

func createPollingMap(sensors []models.Sensor) map[string]int64 {
	pollingMap := make(map[string]int64)

	for _, sensor := range sensors {
		pollingMap[sensor.ID] = 0
	}

	return pollingMap
}
