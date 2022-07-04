package messaging

import (
	"fmt"
	"os"
)

func isRunningInDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

func getMQHost(kind string, user string, pass string) string {
	var host string
	if isRunningInDockerContainer() {
		host = "rabbitmq"
	} else {
		host = "localhost"
	}

	switch kind {
	default:
		return fmt.Sprintf("amqp://%s:%s@%s:5672/", user, pass, host)
	case MQTTKind:
		return fmt.Sprintf("mqtt://%s:%s@%s:1883/", user, pass, host)
	}

}
