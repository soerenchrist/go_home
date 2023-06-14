package background

import (
	"log"
	"time"

	"github.com/soerenchrist/mini_home/db"
	"github.com/soerenchrist/mini_home/models"
)

func PollSensorValues(database db.DevicesDatabase) {
	sensors, err := database.ListPollingSensors()
	if err != nil {
		panic(err)
	}

	lastPolls := createPollingMap(sensors)

	log.Printf("Found %d sensors to poll\n", len(sensors))

	for {
		currentTimestamp := time.Now().Unix()
		for _, sensor := range sensors {
			lastPoll := lastPolls[sensor.ID]
			if currentTimestamp-lastPoll >= int64(sensor.PollingInterval) {
				log.Printf("Polling sensor %s\n", sensor.ID)
				lastPolls[sensor.ID] = currentTimestamp
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func createPollingMap(sensors []models.Sensor) map[string]int64 {
	pollingMap := make(map[string]int64)

	for _, sensor := range sensors {
		pollingMap[sensor.ID] = 0
	}

	return pollingMap
}
