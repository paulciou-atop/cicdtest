package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	//"fmt"
	"io"
	"testing"

	mq "nms/messaging"

	"github.com/stretchr/testify/assert"
)

func write(w io.Writer) chan<- string {
	lines := make(chan string)
	go func() {
		for line := range lines {
			fmt.Fprint(w, line)
		}
	}()
	return lines
}

func Test_AMQPPubSub(t *testing.T) {
	c, err := mq.NewClient() // new client, using amqp protocol as default
	if err != nil {
		t.Error(err)
		return
	}
	defer c.Close() // close connection

	b := new(bytes.Buffer)
	err = c.Subscribe("test", write(b))
	if err != nil {
		t.Error(err)
		return
	}

	msg := "nms test amqp pubsub"
	c.Publish("test", msg)

	// sleep 1s waiting for receving
	time.Sleep(1 * time.Second)
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, string(out), msg, "return is not expected")
}

/*
	This test shows how to receive every topic
*/
func Test_AMQPPubSubWithKey(t *testing.T) {
	c, err := mq.NewClient() // new client, using amqp protocol as default
	if err != nil {
		t.Error(err)
		return
	}
	defer c.Close() // close connection

	b := new(bytes.Buffer)
	err = c.Subscribe("#", write(b)) // receive
	if err != nil {
		t.Error(err)
		return
	}

	msg := "nms test amqp pubsub"
	c.Publish("alan", msg)

	msg2 := "nms test amqp pubsub 2"
	c.Publish("hello", msg2)

	// sleep 1s waiting for receving
	time.Sleep(1 * time.Second)
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
		return
	}

	// only receive msg
	assert.Equal(t, string(out), msg+msg2, "return is not expected")
}

func Test_MQTTPubSub(t *testing.T) {
	c, err := mq.NewClient(mq.MQTTKind) // new mq with protocol mqtt
	if err != nil {
		t.Error(err)
		return
	}
	defer c.Close()

	b := new(bytes.Buffer)
	err = c.Subscribe("nmstest", write(b))
	if err != nil {
		t.Error(err)
		return
	}

	msg := "nms test mqtt pubsub"
	c.Publish("nmstest", msg)

	// sleep 1s waiting for receving
	time.Sleep(1 * time.Second)
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, string(out), msg, "return is not expected")
}
