package messaging

import "errors"

type dummy struct{}

var ErrDummyClient = errors.New("dummy message queue client")

func (d *dummy) Publish(topic string, msg any) error {
	return ErrDummyClient
}

func (d *dummy) Subscribe(topic string, messages chan<- string) error {
	return ErrDummyClient
}

func (d *dummy) Close() error {
	return ErrDummyClient
}

func (d *dummy) GetAMQPClient() (*AMQPClient, error) {
	return nil, errors.New("dummy client")
}
func (d *dummy) GetMQTTClient() (*MQTTClient, error) {
	return nil, errors.New("dummy client")
}
