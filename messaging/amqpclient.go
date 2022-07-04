package messaging

import (
	"context"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Table map[string]interface{}

type ExchangeDeclare struct {
	Name       string
	kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Arguments  Table
}

type QueueDeclare struct {
	Name       string
	Durable    bool
	Exclusive  bool
	AutoDelete bool
	NoWait     bool
	Arguments  Table
}

type QueueBind struct {
	NoWait    bool
	Arguments Table
}

type AMQPPub struct {
	Mandatory bool
	Immediate bool
	Message   amqp.Publishing
}

type AMQPSub struct {
	ConsumerTag string
	NoLocal     bool
	NoAck       bool
	Exclusive   bool
	NoWait      bool
	Arguments   Table
}

type AMQPClient struct {
	Ctx      context.Context
	Pub      AMQPPub
	Sub      AMQPSub
	Exchange ExchangeDeclare
	Queue    QueueDeclare
	Bind     QueueBind

	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (a *AMQPClient) Close() error {
	if a == nil {
		return errors.New("no amqp connection")
	}
	err := a.Channel.Close()
	if err != nil {
		return err
	}

	err = a.Connection.Close()
	if err != nil {
		return err
	}

	return nil
}

func (a *AMQPClient) Publish(topic string, msg interface{}) error {
	if a == nil {
		return errors.New("no amqp connection")
	}
	err := a.Channel.ExchangeDeclare(
		a.Exchange.Name,                  // name
		a.Exchange.kind,                  // type
		a.Exchange.Durable,               // durable
		a.Exchange.AutoDelete,            // auto-deleted
		a.Exchange.Internal,              // internal
		a.Exchange.NoWait,                // no-wait
		amqp.Table(a.Exchange.Arguments), // arguments
	)
	if err != nil {
		return err
	}

	s, ok := msg.(string)
	if !ok {
		err = errors.New("can't convert to string")
		log.Println(err)
		return err
	}

	err = a.Channel.Publish(
		a.Exchange.Name,
		topic,
		a.Pub.Mandatory,
		a.Pub.Immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(s),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *AMQPClient) Subscribe(topic string, messages chan<- string) error {
	if a == nil {
		return errors.New("no amqp connection")
	}

	err := a.Channel.ExchangeDeclare(
		a.Exchange.Name,                  // name
		a.Exchange.kind,                  // type
		a.Exchange.Durable,               // durable
		a.Exchange.AutoDelete,            // auto-deleted
		a.Exchange.Internal,              // internal
		a.Exchange.NoWait,                // no-wait
		amqp.Table(a.Exchange.Arguments), // arguments
	)
	if err != nil {
		return err
	}

	q, err := a.Channel.QueueDeclare(
		a.Queue.Name,
		a.Queue.Durable,
		a.Queue.AutoDelete,
		a.Queue.Exclusive,
		a.Queue.NoWait,
		amqp.Table(a.Queue.Arguments),
	)
	if err != nil {
		log.Println("amqp queue declare failed")
		return err
	}

	err = a.Channel.QueueBind(
		q.Name,          // queue name
		topic,           // routing key
		a.Exchange.Name, // exchange
		a.Bind.NoWait,
		amqp.Table(a.Bind.Arguments),
	)
	if err != nil {
		log.Println("amqp Failed to bind a queue")
		return err
	}

	deliveries, err := a.Channel.Consume(
		q.Name, // queue
		a.Sub.ConsumerTag,
		a.Sub.NoAck,
		a.Sub.Exclusive,
		a.Sub.NoLocal,
		a.Sub.NoWait,
		amqp.Table(a.Sub.Arguments),
	)
	if err != nil {
		log.Println("amqp channel consume failed")
		return err
	}

	go func() {
		for {
			select {
			case <-a.Ctx.Done():
				log.Println("subscribe done")
				close(messages)
				return
			default:
				for d := range deliveries {
					//log.Println("subscribe ", string(d.Body))
					//messages <- d.Body
					messages <- string(d.Body)
				}
			}
		}
	}()

	return nil
}

func connectAmqp(msging *Messaging) error {
	url := getMQHost(AMQPKind, msging.User, msging.Password)

	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	msging.AMQP.Connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	msging.AMQP.Channel = ch

	return nil
}
