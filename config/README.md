## Configuration

Service to allow configuration of ATOP devices.   One at a time, as well as a group of them at a time.  Configuration can be done via APIs.  

There can be CLI command line tools that talk to the API to configure one or more devices.
There can be JSON or YAML files that contain a structured data containing an organized map of configuration data for various devices that express relationships and dependencies.
Also can be called by Web UI that uses REST API.

We will eventually also hook up to other Infrastructure as Code tools such as Ansible, netconf, YANG, terraform, cloudify TOSCA, etc. as needed.

Configuration data and state of current configurations need to be versioned and differences over
different configuration revisions need to be tracked in git or databases.

You should allow for diff'ing of the configurations over time, for example.

We aim to look for ways to detect configuration drift, anomaly and other alarms that require configuration monitoring.


## Build
Build native CLI applicaiton:  
```shell
$ make
```

Build docker image: 
```shell
$ make image
```

## Requirement
- configure a device with specific method
  - configuration with GWD
  - configuration with snmp
  - syslog server settings
  - trap server settings
  - backup device configurations
- configure a collection or a group of devices
  - session management
- schedule configuration
- general API for configuration
  - API `Config.Device` config device with device's ID
  - API `Config.Devices` config group `[device's IDs]`
  - API `Config.Which` config devices which match `[conditions]` (relationship)
    - for example config all devices which `model=switch` and below the subnet `192.168.13.1/24`
    - defined some relation that can describe group fo devices logically
  - API `Config.Store` & `Config.Restore` store/restore configuration with device's ID
  - API `Config.StoreDevices` & `Config.RestoreDevices` store/restore configuration with  group `[device's IDs]`
  - configuration file (YAML or JSON) that describe a collection fo configuration options to group of devices
    - every config services should ignore un-supported options, and return a result of this `wronging`
  - API `Config.Verify` Verify configurations, each config services should have a API allows `configservice` send a config options and check these options available.
  - API `Config.Diff` can diff two configuration 
  - secure (TLS, certificate) and RBAC
- config result with messaging
  - Subscribe topic `config.config` with protocol `amqp`(default) to receive config result

## Command testing
### how to batch config 
`config device` command needs two kind of file, template file and config file. service combines these two file to create configurations. 

config file is a yaml file which describe bunch of device's configuration. In the sample below describes two device and there configurations.

```yaml
- device_id: dev1.json
  device_path: /home/austin/code/atop/NMS/config/pkg/config/dummydevs/dev1.json
  network:
    ip: 10.0.1.11
    mask: 255.255.0.0
  general:
    name: new-dev1
- device_id: dev2.json
  device_path: /home/austin/code/atop/NMS/config/pkg/config/dummydevs/dev2.json
  network:
    ip: 10.0.1.12
  general:
    name: new-dev2
```
template file is a json file consist of configurations, you can get this over `get config` command

```json
go run main.go get config --protocol dummyconfiger --device-id dev --device-path /home/austin/code/atop/NMS/config/pkg/config/dummydevs/dev1.json
{
  "general": {
    "name": "dev1.json"
  },
  "network": {
    "ip": "10.0.12.11",
    "mac": "01:AB:CD:00:1E:10",
    "mask": "255.255.255.0"
  },
  "snmp": {
    "port": 161,
    "private": "private",
    "ver": "V3"
  }
}
```

Let me see how its work:

prerequest:
- dummyconfiger
- postgreSql serivce
- dummy devices, it's very easy to create dummy devices on shell script 
```shell
#!/bin/bash

for i in {1..20};
do
  echo {} > dev$i.json
done
```


- Get template from device
We can get template from devcie  
```shell
$ go run main.go get config -p dummyconfiger -d dev -a /user/asutin/testing/dev1.json > template.json
```

- write a config file, 



Run CLI command
```shell
$ go run main.go config device -c test.yml -t dev1.json
```

## CRUL testing
### Testin with dummyconfiger
Running postgres
```
docker run --rm -d -p 5432:5432 --env POSTGRES_USER=user --env POSTGRES_PASSWORD=pass --env POSTGRES_DB=nms  postgres
```

dummyconfiger supports fields:  
["ip", "name", "mask", "gateway", "snmpEnable"]

Prepare:
dummyconfiger pretends json file as a device. the device.path is file's path, device.id doesn't matter. 


**Config sample:  **
```
$ grpcurl  -d @ -plaintext localhost:8100 devconfig.Config.Device <<EOM
{
  "device":{
    "device_id":"doesn't-matter",
    "device_path":"/home/austin/testing/dev2.json"
  },
  "settings":[
    {
      "protocol":"dummyconfiger",
      "kind":"general",
      "payload":{
        "name":"dev1"
      }
    },
    {
      "kind":"network",
      "payload":{
        "ip":"localhost",
        "mask":"255.255.255.0",
        "gateway":"10.10.1.254"
      }
    }
  ]
}
EOM

{
  "session": {
    "id": "ab0caaac-549c-40d2-875d-86bcce68ef3c",
    "state": "running",
    "startedTime": "2022-06-05 01:21:39"
  }
}
```

**Config fail sample:**
```
$ grpcurl  -d @ -plaintext localhost:8100 devconfig.Config.Device <<EOM
{
  "device":{
    "device_id":"doesn't-matter",
    "device_path":"/Users/austinjan/testing/dev1.json"
  },
  "settings":[
    {
      "protocol":"dummyconfiger",
      "kind":"general",
      "payload":{
        "name":"dev1",
        "unsupported":"whatever"
      }
    },
    {
      "kind":"network",
      "payload":{
        "ip":"localhost",
        "mask":"255.255.255.0",
        "gateway":"10.10.1.254"
      }
    }
  ]
}
EOM
```

**Get result**
```shell
$ grpcurl -d '{"session_id":"57d3a2f7-d822-47b4-8430-11fc4761b707"}' -plaintext localhost:8100 devconfig.Config.GetResult
```

**Get session list**
```shell
grpcurl -d '{}' -plaintext localhost:8100 devconfig.Config.List
```

Config devices
```shell
$ grpcurl  -d @ -plaintext localhost:8100 devconfig.Config.Devices <<EOM
{
  "devices":[
    {
      "device_id":"doesn't-matter",
      "device_path":"/home/austin/testing/dev1.json"
    },
    {
      "device_id":"doesn't-matter",
      "device_path":"/home/austin/testing/dev3.json"
    }
  ],
  "settings":[
    {
      "protocol":"dummyconfiger",
      "kind":"general",
      "payload":{
        "name":"dev1-multiple"
      }
    }
  ]
}
EOM


{
  "device": {
    "deviceId": "doesn't-matter",
    "devicePath": "/home/austin/testing/dev1.json"
  },
  "session": {
    "id": "eee40f9f-4afe-4527-98a8-6c8381c54a1c",
    "state": "running",
    "startedTime": "2022-06-06 17:59:51"
  }
}
{
  "device": {
    "deviceId": "doesn't-matter",
    "devicePath": "/home/austin/testing/dev2.json"
  },
  "session": {
    "id": "f45b302a-20ec-4ad4-9387-e7704ad368e6",
    "state": "running",
    "startedTime": "2022-06-06 17:59:51"
  }
}

```

## Subscribe Config

```
$ grpcurl  -d @ -plaintext localhost:8100 devconfig.Config.Device <<EOM
{
  "device":{
    "device_id":"doesn't-matter",
    "device_path":"/home/austin/testing/dev2.json"
  },
  "settings":[
    {
      "protocol":"dummyconfiger",
      "kind":"general",
      "payload":{
EOM } } "gateway":"10.10.1.254"
{
  "device": {
    "deviceId": "doesn't-matter",
    "devicePath": "/home/austin/testing/dev2.json"
  },
  "session": {
    "id": "40529b24-c722-4a04-a963-09befe55eb68",
    "state": "running",
    "startedTime": "2022-06-09 10:35:32"
  }
}
```

```
$ go run messaging/main.go sub -t config.config
Starting to receive from amqp topic config.config
Press Ctrl + C to stop receiving...
{
  "SessionID": "40529b24-c722-4a04-a963-09befe55eb68",
  "State": "fail",
  "StartedAt": "2022-06-09 10:35:32",
  "EndedAt": "2022-06-09 10:35:32",
  "Message": "can not reach device /home/austin/testing/dev2.json"
}
```
