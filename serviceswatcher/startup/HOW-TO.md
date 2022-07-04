# Integration Guide

This document defined a convention of each service to make sure services are able to work with each others later.

There are some convention services have to follow:

- Every service has to have a Dockerfile. And put it in the service's root folder.
- Every service has to have `run` command, otherwise author have to modify `docker-compose.yaml` file to expose how to start up the service.
- Every service has to have a `.proto` file which describe the APIs information. Put it in the `/NMS/src/services/api/proto/{version}`

## Step of communicate with other services

1. Connect to `serviceswatcher` service

   ```go
   conn, err := grpc.Dial("serviceswatcher:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
   if err != nil {
     glog.Error("Can not connect to serviceswatcher")
     return err
   }
   defer conn.Close()
   client := NewWatcherClient(conn)
   ```

2. Get specifc service's host informaton

   ```go
   res,err:= client.Get(ctx,&GetRequest{
   ServiceName: "snmpscan",
   })
   
   if err != nil {
     glog.Error("Call servicewatcher.Watcher.Get err: ", err)
     return err
   }
   
   host := fmt.Sprintf("%s:%d", req.Info.Address, req.Info.Port)
   ```
3. Create connection with service, and create client to access API.

## Dockerfile
- Every service should have a workable Dockerfile in the root folder
- Default hostname will be service name which is service's folder name in the `/NMS/src/services`
- If service need to communicate with network out side container, you can modify `/NMS/src/services/servicewatcher/startup/docker-compose.yml` file. Add `port` field in your service section.
- Every service's default listen port should be `8080`
  
## Finding servicewatcher
By default servicewathcer's hostname is `servicewathcer` and default port is `8081`, please refer a sample code below to figure out how to get the `servicewatcher`

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

## Database
There is a Postgresql in docker compose. You can use this service directly when docker is up.

The Postgresql settings can refer to docker-compose.yaml

## Rebuild
You should rebuild container every time when you update source code. Running command below on this folder (`/NMS/src/services/serviceswatcher/startup`)

`docker-compose build`

## Startup 
`docker-compose up`