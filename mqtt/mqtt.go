package mqtt

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttController struct {
	publishChannel PublishChannel
}

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

func ConnectToBroker(config MqttConfig, publish PublishChannel) error {
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
