dummyconfiger is a fake congier service, this service do nothing but write config metric from `config` service to log file

This service take a json file as a device, config device = update file's content.

## Build
If you want to build this service in the windows os, add postfix `.exe` on the output file name

In the path `/nms/testing/dummyconfiger`  
```shell
$ go build -o ./bin/dummyconfiger ./cmd/dummyconfiger
```
## Run
In the path `/nms/testing/dummyconfiger`  
```shell
$ ./bin/dummyconfiger
```

## grpcrul usage sample
```shell
$ grpcurl -d '{"session":{"state":"running","id":"123","startedTime":"2022/6/1"},"device":{"device_id":"dev1","device_path":"test.json"},"configs":[{"protocol":"file","kind":"testing","payload":{"name":"dummy"}}]}' -plaintext localhost:8085 configer.Configer.Config

{
  "session": {
    "id": "123",
    "state": "success",
    "startedTime": "2022/6/1"
  }
}

```

```shell
$ grpcurl -d '{"session":{"state":"running","id":"123","startedTime":"2022/6/1"},"configs":[{"protocol":"file","kind":"testing","payload":{"name":"dummy","age":17}}]}' -plaintext localhost:8085 configer.Configer.Validate
{
  "session": {
    "id": "123",
    "state": "success",
    "startedTime": "2022/6/1"
  },
  "failConfigs": [
    {
      "protocol": "file",
      "kind": "testing",
      "failFields": [
        "age"
      ]
    }
  ]
}

```

```shell
$  grpcurl -d '{"device":{"device_id":"123", "device_path":"test.json"},"kinds":["testing"]}' -plaint
ext localhost:8085 configer.Configer.GetConfig
{
  "device": {
    "deviceId": "123",
    "devicePath": "test.json"
  },
  "configs": {
      "testing": {
            "ip": "192.10.10.1",
            "name": "dummy"
          }
    }
}


```