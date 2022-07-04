services service to keep each services's basic information & status

- List each service informaiton



## TODO
- heartbeat API 
- store data into datastore service

## Strategy
- `servicesswathcer` running on fixed port (8081) and fixed host name `serviceswatcher`
- Every service could register/update itself with Register API
- Every service could query any service with 

## Folders & files
**Folders**
- /api: gRPC APIs and clients
- /cmd: CLI command main file
- /config: viper & service default configuration
- /pkg: service's logic package, functions etc...
- /startup: docker-compose config file to deploy NMS and testing service
- /script: building script
- /prometheus: prometheus stuff

**Files**
- Makefile: make gRPC API, build container or run container
- README.md: this file

## Build project
**Build**
Build native CLI applicaiton:  
```shell
$ make
```

Build docker image: 
```shell
$ make image
```

**Run**
**Native**
Binary executable file was output on `.\bin` sub-folder:
```
$ .\bin\sw
```

## Running service
1. Build docker image
  In the root path of `/serivceswatcher` running command to build image.
  ```shell
  $ docker build --no-cache -t serviceswatcher:v1 .
  ```
2. Running container
   Because we want `serviceswatcher` always has a fixed host name let other services connect, it's better give it a fixed hostname with customized network. should be 172.18.0.2:8081
   ```shell
   # create network
   $ docker network create --subnet=172.18.0.0/16 demosite
   $ docker run --rm --net demosite -h serviceswatcher --ip 172.18.0.2 -p 8081:8081 serviceswatcher:v1 run
   ```
   > Testing network settings is corret.
   ```shell
   $ docker run --rm --net -it demosite ubuntu /bin/bash
   # ping with hostname
   root@6b80981c64c6:> apt-get upgrade
   root@6b80981c64c6:> apt-get install iputils-ping
   root@6b80981c64c6:> ping serviceswatcher
    # As you can see, serviceswatcher's IP is 172.18.0.2
    PING serviceswatcher (172.18.0.2) 56(84) bytes of data.
    64 bytes from 253f19c88a5f.demosite (172.18.0.2): icmp_seq=1 ttl=64 time=0.075 ms
   # traceroute with hostname
   root@6b80981c64c6:> apt-get install traceroute
   root@6b80981c64c6:> traceroute serviceswatcher
   traceroute to serviceswatcher (172.18.0.2), 30 hops max, 60 byte packets
     1  253f19c88a5f.demosite (172.18.0.2)  0.127 ms  0.099 ms  0.047 ms
   ```

## Setup default services information
We want to provide default services information, please add services information in `settings.json` of this service. For milestone-1 scenario we will running every services as a container, its easily to communicate each other when we assign a fixed host name for each services and run services with same network. Let's see how to fulfill it.

### Running service with specific host name

Running snmpscan service witch host name is `snmpscan` with specific network (demosite)
```shell
$ docker run --rm --net demosite -h snmpscan  snmpscan:v1 run
```

### Add service information into `settings.json`
After start up your service as a container with specific host name, you have to add this informaiton into `/servicewatcher/setting.json`, for example `snmpscan` service should add information into `setting.json` like this:  

```json
{
  "services": [
    {
      "name": "servicewatcher",
      "address": "servicewatcher",
      "port": 8081,
      "kind": [
        "https",
        "grpc"
      ]
    },
    {
      "name": "snmpscan",
      "address": "snmpscan",
      "port": 8088,
      "kind": [
        "grpc","http"
      ]
    }
  ]
}
```  

## Testing scenario

1. Run `serviceswatcher`
```shell
$ docker run --rm --net demosite -h snmpscan  snmpscan:v1 run
```

2. Execute register command like  `snmpscan register --address=snmpscan --port=1234` to change service's informaiton in the `servicesswatcher`, keep in mind set same network when you run container.  
```shell
$  docker run --rm --net demosite -h snmpscan  snmpscan:v1 register --address=snmpscan --port=1234
service [snmpscan] was registered on snmpscan:1234
```

3. Use grpcurl to check service's information was changed
Start grpcurl container
```shell
# Keep in mind grpcrul must have same network setting
$ docker run --rm -it --net=demosite networld/grpcurl
# Into container terminel, use grpcurl to check result
> ./grpcurl -plaintext serviceswatcher:8081 serviceswatcher.Watcher.List
{
  "infos": [
    {
      "name": "servicewatcher",
      "address": "servicewatcher",
      "port": 8081,
      "kind": [
        "https",
        "grpc"
      ]
    },
    {
      "name": "snmpscan",
      "address": "snmpscan",
      "port": 1234,
      "kind": [
        "grpc"
      ]
    }
  ]
}
```
## Samples
### Sample of call services which has registered in the serviceswatcher
1. Compile `serviceswatcher.proto` , you should change some path setting to suit your project.
```
	protoc -I ../api/proto --go_out . --go-grpc_out . \
		--grpc-gateway_out . \
		--openapiv2_out ./api/doc --openapiv2_opt allow_merge=true,merge_file_name=api \
		../api/proto/v1/serviceswatcher.proto
```

2. Implement client code to finding sepcific service
```go
func Store(data map[string]interface{}) error {

  // Get service info from serviceswatcher
	conn, err := grpc.Dial("serviceswatcher:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
			glog.Error("Can not connect to serviceswatcher")
			return err
	}
	defer conn.Close()
	client := NewWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
  // replace ServiceName to which you want to investigate
	res,err:= client.Get(ctx,&GetRequest{
	ServiceName: "snmpscan",
  })

	if err != nil {
		glog.Error("Call servicewatcher.Watcher.Get err: ", err)
		return err
	}

  host := fmt.Sprintf("%s:%d", req.Info.Address, req.Info.Port)

  // connect to service 
  conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
			glog.Errorf("Can not connect to %s",host)
			return err
	}
	defer conn.Close()
	// ... place your code here
  // ...

  return nil
}
```

## Sample for develop
When you develop service which will ask `serviceswatcher`, you should not running all services every time. In this chapter we describe how to enhance your develop iteration flow.

**Pre-requirement**
- create a network
 ```shell
 $ docker network create testing-net
 ```

**Step**
1. Running `servicewatcher`
   - Build image from source code, in the `/NMS/src/services/serviceswatcher` execute command:
   ```shell
   $ docker build -t sw .
   ```
   - Running `serviceswatcher`, execute command: 
   ```shell
   $ docker run -d --rm --network testing-net -h serviceswatcher sw run
   ```
2. Build your service and repeat steps 1.

## Put every thing together
In the folder `./startup` 

