package mqtt

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
	"github.com/soerenchrist/go_home/util"
)

type MqttConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	ClientId string
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

type Message struct {
	Topic   string
	Payload string
}

type PublishChannel chan Message

func AddMqttBinding(config MqttConfig, publish PublishChannel, database db.Database, outputBindings chan models.SensorValue) error {
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("ssl://%s:%d", config.Host, config.Port))
	options.SetClientID(config.ClientId)
	options.SetUsername(config.Username)
	options.SetPassword(config.Password)
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(options)

	log.Printf("Connecting to MQTT broker at %s:%d", config.Host, config.Port)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Println("Connected to MQTT broker... Listening for publishes")

	go listenForPublishes(client, publish)
	subscribeToNewValues(client, database, outputBindings)

	return nil
}

func listenForPublishes(client mqtt.Client, publish PublishChannel) {
	for {
		message := <-publish
		log.Println("Publishing message: ", message)
		token := client.Publish(message.Topic, 0, false, message.Payload)
		token.Wait()
	}
}

func subscribeToNewValues(client mqtt.Client, database db.Database, outputBindings chan models.SensorValue) {
	client.Subscribe("home/+/+/value", 0, func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		parts := strings.Split(topic, "/")
		deviceId := parts[1]
		sensorId := parts[2]

		value := string(msg.Payload())
		log.Printf("Received value %s for sensor %s on device %s", value, sensorId, deviceId)

		sensor, err := database.GetSensor(deviceId, sensorId)
		if err != nil {
			log.Println("Error getting sensor: ", err)
			return
		}

		sensorValue := &models.SensorValue{
			SensorID:  sensor.ID,
			DeviceID:  sensor.DeviceID,
			Value:     value,
			Timestamp: util.GetTimestamp(),
		}

		switch sensor.DataType {
		case models.DataTypeFloat:
			_, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Println("Error parsing float: ", err)
				return
			}
			database.AddSensorValue(sensorValue)
		case models.DataTypeInt:
			_, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Println("Error parsing int: ", err)
				return
			}
			database.AddSensorValue(sensorValue)
		case models.DataTypeBool:
			_, err := strconv.ParseBool(value)
			if err != nil {
				log.Println("Error parsing bool: ", err)
				return
			}
			database.AddSensorValue(sensorValue)
		default:
			database.AddSensorValue(sensorValue)
		}
		outputBindings <- *sensorValue
	})
}
