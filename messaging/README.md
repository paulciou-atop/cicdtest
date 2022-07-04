# Support for Messaging protocols and APIs

MQTT, AMQP and other types of messaging or pub/sub based communication code that reside under
a neutral APIs to make them appear generic.


## NMS Current Message Queue

### AMQP

| Service  |      Topic           |     Message    |
|:---------|:--------------------:|----------------|
| config   |  config.config       | config result (json)|
| snmpscan |    scan.snmpscan     | {"sessionid":"xxxx"}  (json)|
| atopudpscan | scan.atopudpscan  | {"sessionid":"xxxx"}  (json)|

Please refer to each service README for more Message details

- Config
	Subscribe topic `config.config` with protocol `amqp`(default) to receive config result

- Scan
	- SnmpScan 

		Subscribe topic `scan.snmpscan` with protocol `amqp`(default) to receive message when snmpscan service scan finished
	- AtopUdpScan

		Subscribe topic `scan.atopudpscan` with protocol `amqp`(default) to receive message when atopudpscan service scan finished
	- All

		Subscribe topic `scan.*` with protocol `amqp`(default) to receive one level hierarchy topic or `scan.#` to receive all topic after `scan.`.
		Please refer to [RabbitMQ Topic](https://www.rabbitmq.com/tutorials/tutorial-five-go.html) for more details

### MQTT

None


## How to use
Messaging is a package which can be imported and used directly
Refer to the test code in `test/messaging_test.go`

### New Client
create a new messaging client with deafult settings
```go
c, err := mq.NewClient()
```

This will create a client with amqp kind by default, if you want to create a mqtt one's, just assign `messaging.MQTTKind` to the NewClient function.

### Publish

Publish is easy to use, provide `topic` and `message` and it's good to go. 


### Subscribe

Subscribe requires not only a topic but also a channel to receive the messages continuously. Please refer to test or cmd code.



### Topic Naming rules


NMS AMQP messaging uses dot to seperate the topic hierarchy for example scan.snmpscan, scan.atopudpscan, please follow the rule
`"purpose"."service name"` to name a topic

On the other hand, MQTT uses "/" as a delimiter, please follow `"purpose"/"service name"` rule, for example
scan/snmpscan, scan/atopudpscan, config/config .. etc



## Testing

### Run docker
There is a docker-comnpose.yml file in test folder, run command below first to start the rabbitmq server

```shell
docker-compose up
```

### Test function
In the test folder run

```shell
go test
```

### Command line interface
Messaging also provides CLI in the cmd folder run 

```shell
go run messaging/main.go
```
or make or build.bat to generate the excutive file depends on the develop environment

There will be two commands, pub and sub, these two commands can used for publish and subscribe through rabbitmq server in docker with specific protocol and topic.


## TODO

- Currently most parameters are set there like user, password, host ...etc. Planning to load configurations from file or config service
- Refactoring the code ,interface and data structure
