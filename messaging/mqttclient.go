package messaging

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTPub struct {
	Qos     byte
	Retaind bool
}

type MQTTSub struct {
	Qos byte
}

type DisconnectOpt struct {
	Quiesce uint
}

type MQTTClient struct {
	Pub        MQTTPub
	Sub        MQTTSub
	Disconnect DisconnectOpt
	Connection mqtt.Client
}

func (m *MQTTClient) Close() error {
	if m == nil {
		return errors.New("no mqtt client")
	}
	m.Connection.Disconnect(m.Disconnect.Quiesce)
	return nil
}

func (m *MQTTClient) Publish(topic string, msg interface{}) error {
	if m == nil {
		return errors.New("no mqtt client")
	}
	token := m.Connection.Publish(topic, m.Pub.Qos, m.Pub.Retaind, msg)
	for !token.WaitTimeout(time.Microsecond) {
	}
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MQTTClient) Subscribe(topic string, messages chan<- string) error {
	if m == nil {
		return errors.New("no mqtt client")
	}
	token := m.Connection.Subscribe(
		topic,
		m.Sub.Qos,
		func(messages chan<- string) func(mqtt.Client, mqtt.Message) {
			return func(c mqtt.Client, m mqtt.Message) {
				//messages <- m.Payload()
				messages <- string(m.Payload())
			}
		}(messages),
	)
	for !token.WaitTimeout(time.Microsecond) {
	}
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

func connectMqtt(msging *Messaging) error {
	mqttURL := getMQHost(MQTTKind, msging.User, msging.Password)
	uri, err := url.Parse(mqttURL)
	if err != nil {
		return err
	}

	password, _ := uri.User.Password()
	name := uri.User.Username()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(name)
	opts.SetPassword(password)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(time.Microsecond) {
	}
	if token.Error() != nil {
		return token.Error()
	}

	msging.MQTT.Connection = client
	return nil
}
