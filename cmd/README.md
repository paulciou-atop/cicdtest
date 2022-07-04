## Various command line utilities
The CLI tools that use APIs to perform tasks that are available via API endpoints.

The same APIs can be used by the Web apps. Both CLI and web clients have the same APIs.

This allows for implementing APIs that will be accessible from CLI as well as Web UI and also ease the testing and verification process.


### Pre requirement
Go version : >1.18  

### Build CLI

```shell
$ sh build.sh
```

#### Compile proto
Please refer Makefile 

```shell
$ make grpc
```

### Connect Services

#### Native 
Has to setup the ports correctly and register services to serviceswatcher in order to connect each other

```shell
$ src/services/serviceswatcher go run main.go run
$ src/services/snmpscan go run main.go run
$ src/services/scanservice go run main.go run
```

#### Dockerized
Start docker compose in serviceswatcher/startup, please refer to serviceswatcher/startup/HOW-TO for more details

```shell
$ docker-compose build
$ docker-compose up
```

### CLI

#### Scan Service

```shell
$ nmsctl scan start [service-address] [snmp-range] [gwd-server-ip]
$ nmsctl scan status [service-address] [sessionid]
$ nmsctl scan stop [service-address] [sessionid]
$ nmsctl scan result [service-address] [sessionid]
```

#### SNMP Scan

```shell
$ nmsctl snmpscan get [service-adress] [target]
$ nmsctl snmpscan scan [service-adress] [range]
$ nmsctl snmpscan walk [service-adress] [target] [root oid]
```

#### UDP Scan

```shell
$ nmsctl udp beep [service-adress] [device ip] [device mac] [server ip]
$ nmsctl udp ip [service-adress]
$ nmsctl udp scan [service-adress] [server ip] [timeout]
$ nmsctl udp stop [service-adress]

```

#### Services Watcher

```shell
$ nmsctl service list [service-adress]
$ nmsctl service status [service-adress] [name]
```

### Test CLI with snmp simulators
Go to serviceswatcher/startup

add snmpsim in docker-compose.yml
```
  snmpsim:
    build: ../../snmpscan/test/testdata
```

```shell
$ docker-compose up --scale snmpsim=10
```

This will run the services and ten SNMP simulators.

And then we can use scan command to find them with range 172.19.0.0/24 which is the docker compose default network
