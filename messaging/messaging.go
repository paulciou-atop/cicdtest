package messaging

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrNullMQ       = errors.New("got null mq")
	ErrMQNotSupport = errors.New("message queue not support")
)

const (
	AMQPKind = "amqp"
	MQTTKind = "mqtt"
)

type Publisher interface {
	Publish(topic string, msg interface{}) error
}

type Subscriber interface {
	Subscribe(topic string, messages chan<- string) error
}

type IMQAccess interface {
	Publisher
	Subscriber
}

type IMQClient interface {
	IMQAccess
	GetAMQPClient() (*AMQPClient, error)
	GetMQTTClient() (*MQTTClient, error)
	Close() error
}

type Messaging struct {
	Host     string
	Port     uint
	User     string
	Password string
	Kind     string

	AMQP AMQPClient
	MQTT MQTTClient
}

var defaultMessaging = Messaging{
	User:     "guest",
	Password: "guest",
	Kind:     AMQPKind,
	AMQP: AMQPClient{
		Ctx: context.Background(),
		Pub: AMQPPub{
			Mandatory: false,
			Immediate: false,
		},
		Sub: AMQPSub{
			ConsumerTag: "",
			NoLocal:     false,
			NoAck:       true,
			Exclusive:   false,
			NoWait:      false,
			Arguments:   nil,
		},
		Exchange: ExchangeDeclare{
			Name:       "nms",   // default name
			kind:       "topic", // "direct", "fanout", "topic" and "headers", can't be modified after exchange declared
			Durable:    true,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Arguments:  nil,
		},
		Queue: QueueDeclare{
			Name:       "",
			Durable:    false,
			Exclusive:  true,
			AutoDelete: false,
			NoWait:     false,
			Arguments:  nil,
		},
		Bind: QueueBind{
			NoWait:    false,
			Arguments: nil,
		},
	},
	// Mqtt: MqttOptions{

	// },
}

func (m *Messaging) Close() error {
	switch m.Kind {
	case AMQPKind:
		return m.AMQP.Close()
	case MQTTKind:
		return m.MQTT.Close()
	}

	return ErrMQNotSupport
}

func (m *Messaging) Publish(topic string, msg interface{}) error {
	switch m.Kind {
	case AMQPKind:
		return m.AMQP.Publish(topic, msg)
	case MQTTKind:
		return m.MQTT.Publish(topic, msg)
	}

	return ErrMQNotSupport
}

func (m *Messaging) Subscribe(topic string, messages chan<- string) error {
	switch m.Kind {
	case AMQPKind:
		return m.AMQP.Subscribe(topic, messages)
	case MQTTKind:
		return m.MQTT.Subscribe(topic, messages)
	}

	return ErrMQNotSupport
}

func (m *Messaging) GetAMQPClient() (*AMQPClient, error) {
	return &m.AMQP, nil
}

func (m *Messaging) GetMQTTClient() (*MQTTClient, error) {
	return &m.MQTT, nil
}

// connect to Rabbitmq immediately and direclty
func NewClient(kind ...string) (IMQClient, error) {
	o := defaultMessaging
	dmy := &dummy{}
	if len(kind) >= 1 {
		o.Kind = kind[0]
	} else {
		o.Kind = AMQPKind
	}

	switch o.Kind {
	case AMQPKind:
		err := connectAmqp(&o)
		if err != nil {
			return dmy, err
		}

		return &o, nil
	case MQTTKind:
		err := connectMqtt(&o)
		if err != nil {
			return dmy, err
		}
		return &o, nil
	}
	return dmy, fmt.Errorf("mq: %s create client fail: %w", o.Kind, ErrMQNotSupport)
}

var retry = 3

// provide a context with timeout to dial to Rabbitmq in retry times
// return connection on success or dummy
func DialMessaging(ctx context.Context, kind ...string) IMQClient {
	c := make(chan IMQClient)
	var k string
	if len(kind) >= 1 {
		k = kind[0]
	} else {
		k = AMQPKind
	}

	go func() {
		for {
			client, err := NewClient(k)
			if err == nil {
				c <- client
				return
			}
			time.Sleep(time.Second * 3)
		}
	}()
	select {
	case <-ctx.Done():
		log.Info("cancel connect to Rabbitmq, use dummy client instead")
		close(c)
		return &dummy{}

	case client := <-c:
		close(c)
		return client
	}

}
