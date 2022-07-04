## ATOP UDP Scan service

Send out UDP broadcast to find ATOP devices on the directly attached subnet



## Make

```makefile
Usage:
"make  <target>"
Targets:
bin/atopudpscan: build bin/atopudpscan
build:          build atopudpscan
run:            run atopudpscan service 
image:  docker build atopudpscan
docker-run:     docker run atopudpscan
test:           run test atopudpscan service
clean:          clean test data
grpc:           build *.proto file
```



## GRPC api

#### GwdClient

```go
type GwdClient interface {
    //make device sound
	Beep(ctx context.Context, in *GwdConfig, opts ...grpc.CallOption) (*Response, error)
    //return Server Ip
	GetServerIp(ctx context.Context, in *EmptyParams, opts ...grpc.CallOption) (*ServerIp, error)
    //Stop Session scan
	Stop(ctx context.Context, in *Sessions, opts ...grpc.CallOption) (*Response, error)
    //session scan
	SessionScan(ctx context.Context, in *ScanConfig, opts ...grpc.CallOption) (*ResponseSession, error)
    //Get SessionStatus
	GetSessionStatus(ctx context.Context, in *Sessions, opts ...grpc.CallOption) (*ResponseSession, error)
    //GetSessiondata from datastore
	GetSessionData(ctx context.Context, in *Sessions, opts ...grpc.CallOption) (*DeviceResponse, error)
}

```



#### AtopDeviceClient

```go
type AtopDeviceClient interface {
    //set device config
	SettingConfig(ctx context.Context, in *GwdConfig, opts ...grpc.CallOption) (*Response, error)
    //upload file and return new filename
	Upload(ctx context.Context, opts ...grpc.CallOption) (AtopDevice_UploadClient, error)
    //device fw upgrade after complete file deleted
	FwUpgrading(ctx context.Context, in *FwInfo, opts ...grpc.CallOption) (*Response, error)
    //get upgrading status of device
	GetProcessStatus(ctx context.Context, in *FwRequest, opts ...grpc.CallOption) (*FwMessage, error)
    //reboot device
	Reboot(ctx context.Context, in *GwdConfig, opts ...grpc.CallOption) (*Response, error)
    //get server ips
	GetServerIp(ctx context.Context, in *EmptyParams, opts ...grpc.CallOption) (*ServerIp, error)
}
```


## Messaging

Atopudpscan publish messages while scanning is done no matter successed or failed

- Subscribe topic `scan.atopudpscan` with protocol `amqp`(default) to receive message when atopudpscan service scan finished

### Pre requirements

- Rabbitmq docker service running

### Example

```json
{
  "sessionid": "ss:b37f607a-5f0f-4076-8ed2-6239b8c7650b",
}
```