package mqtt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
	"github.com/soerenchrist/go_home/internal/value"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("mqtt")

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
	log.Debugf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Debugf("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Debugf("Connect lost: %v", err)
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

	log.Infof("Connecting to MQTT broker at %s:%d", mqttConf.Host, mqttConf.Port)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Info("Connected to MQTT broker... Listening for publishes")

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
			log.Errorf("Error while marshalling request: %s", err)
		}

		host := config.GetString("server.host")
		port := config.GetInt("server.port")

		url := fmt.Sprintf("http://%s:%d/api/v1/devices/%s/sensors/%s/values", host, port, device, sensor)

		_, err = http.Post(url, "application/json", strings.NewReader(string(body)))
		if err != nil {
			log.Errorf("Error while sending request to server: %s", err)
		}
	})
	token.Wait()
	log.Debugf("Subscribed to topic: %s\n", topic)

}

func listenForPublishes(client mqtt.Client, publish PublishChannel) {
	for {
		message := <-publish
		log.Debug("Publishing message: ", message)
		token := client.Publish(message.Topic, 0, false, message.Payload)
		token.Wait()
	}
}
