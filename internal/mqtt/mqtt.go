package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/soerenchrist/go_home/internal/value"
	"github.com/spf13/viper"
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

func ConnectToBroker(mqttConf MqttConfig, publish PublishChannel, config *viper.Viper) error {
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("ssl://%s:%d", mqttConf.Host, mqttConf.Port))
	options.SetClientID(mqttConf.ClientId)
	options.SetUsername(mqttConf.Username)
	options.SetPassword(mqttConf.Password)
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(options)

	log.Printf("Connecting to MQTT broker at %s:%d", mqttConf.Host, mqttConf.Port)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Println("Connected to MQTT broker... Listening for publishes")

	go listenForPublishes(client, publish)
	subscribe(client, config)

	return nil
}

func subscribe(client mqtt.Client, config *viper.Viper) {
	topic := "home/+/+/data"
	token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		parts := strings.Split(msg.Topic(), "/")
		device := parts[1]
		sensor := parts[2]

		request := value.AddSensorValueRequest{
			Value: string(msg.Payload()),
		}
		body, err := json.Marshal(request)
		if err != nil {
			log.Printf("Error while marshalling request: %s", err)
		}

		host := config.GetString("server.host")
		port := config.GetInt("server.port")

		url := fmt.Sprintf("http://%s:%d/api/v1/devices/%s/sensors/%s/values", host, port, device, sensor)

		_, err = http.Post(url, "application/json", strings.NewReader(string(body)))
		if err != nil {
			log.Printf("Error while sending request to server: %s", err)
		}
	})
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)

}

func listenForPublishes(client mqtt.Client, publish PublishChannel) {
	for {
		message := <-publish
		log.Println("Publishing message: ", message)
		token := client.Publish(message.Topic, 0, false, message.Payload)
		token.Wait()
	}
}
