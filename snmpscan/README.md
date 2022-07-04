## SNMP Scanner service

Uses SNMP to find devices on the network.  

### Pre requirement
Go version : >1.18  

### TODO
- secure gRPC
- secure HTTP
  
### Folders
- /test: testing code here  
- /api: gRPC API code  
  - /grpcsnmp: SNMP APIs  
  - /server: http & gRPC server code  
- /cmd: CLI stuffs  
- /pkg:  packages here

### Build container
In the /snmpscan:  

1. Build a container
```shell
$ docker build --tag snmpscan:test .
```

Run snmpscan service
> If you want to test `snmpscan` with other services you should attach suitable network, please refer 'serviceswather' service's READMD.md for more information
```shell
$ docker run --rm --network={network} -h snmpscan snmpscan:test run
```

Excute CLI command get, to get 192.168.2.1's SNMP information  
```shell
$ docker run --rm snmpscan:test get snmp 192.168.2.1
```

Scan atop device with CIDR notation
```shell
$ docker run --rm snmpscan:latest scan --atop-device --range=192.168.13.1/24
```
Async scan for atop devices, please use your custom network instead of `testing-net` in the sample below.
```shell
# running redis first
$ docker run --network={network-name} -d -h redis redis
# running serviceswatcher
$ docker run --network={network-name} -d -h serviceswatcher serviceswatcher:latest
# running async scan
$ docker run --network=testing-net test scan --async --atop-devices --range=192.168.13.1/24
# get result
$ docker run --network=testing-net test get session {session-id}
```

### Compile proto
Please refer Makefile 

```shell
$ make grpc
```

In makefile every thing will output on /api/{service}/v1 folder, which was defined in .proto file `option go_package`


### Messaging

Snmpscan publish messages while scanning is done

- Subscribe topic `scan.snmpscan` with protocol `amqp`(default) to receive message when snmpscan service scan finished

#### Pre requirements

- Rabbitmq docker service running

#### Example

```
$ grpcurl -d '{"range":"192.168.13.221/24","atop_devices":true}' -plaintext localhost:8084 snmpscan.SnmpScan.StartAsyncScan
{
  "sessionId": "ss:b37f607a-5f0f-4076-8ed2-6239b8c7650b",
  "success": true
}
```

```
$ go run messaging/cmd/messaging/main.go sub -t scan.snmpscan
Starting to receive from amqp topic scan.snmpscan
Press Ctrl + C to stop receiving...
{
  "sessionid": "ss:b37f607a-5f0f-4076-8ed2-6239b8c7650b",
}
```


### Testing CLI

#### 1. Build snmpsim container
Go into ./test/testdata/

Build image and assign name = snmpsim:v1 
```shell
$ docker build -t snmpsim:v1 .
```
> If you want to change SNMP agent response, you can modify public.snmprec 

Run image you build previously. You can run container as many as you want, keep in mind if you want to run multiple containers, you should not give container name with `--name`
```shell
 $ docker run -d --rm --name=snmpsim snmpsim:v1 
```

Option: Check IP Address (Should be 172.17.0.x).
```shell
$ docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' snmpsim
```

Scan SNMP in whole container sub-net.
```shell
$ docker run --rm {snmpscan-image} scan --range=172.17.0.1/24
```

You should see result like this:
```json
[
  [
    {
      "value": "HP ETHERNET MULTI-ENVIRONMENT,SN:VNBNM37B77,FN:5L5GJTL,SVCID:10288,PID:HP Color LaserJet MFP M281fdw",
      "name": ".1.3.6.1.2.1.1.1.0",
      "kind": "string",
      "oid": ".1.3.6.1.2.1.1.1.0"
    },
    {
      "value": ".1.3.6.1.4.1.11.2.3.9.1",
      "name": ".1.3.6.1.2.1.1.2.0",
      "kind": "string",
      "oid": ".1.3.6.1.2.1.1.2.0"
    },
    {
      "value": 158180676,
      "name": ".1.3.6.1.2.1.1.3.0",
      "kind": "int",
      "oid": ".1.3.6.1.2.1.1.3.0"
    },
    {
      "value": "",
      "name": ".1.3.6.1.2.1.1.4.0",
      "kind": "string",
      "oid": ".1.3.6.1.2.1.1.4.0"
    },
    {
      "value": "NPI405F13",
      "name": ".1.3.6.1.2.1.1.5.0",
      "kind": "string",
      "oid": ".1.3.6.1.2.1.1.5.0"
    },
    {
      "value": "C4-65-16-40-5F-13",
      "name": "MAC",
      "kind": "string",
      "oid": ""
    },
    {
      "value": "192.168.13.200",
      "name": "IP",
      "kind": "string",
      "oid": ""
    }
  ],
  [
    {
      "value": "Managed Switch, EHG7508-8PoE",
      "name": ".1.3.6.1.2.1.1.1.0",
      "kind": "string",
      "oid": ".1.3.6.1.2.1.1.1.0"
    },
    {
      "value": ".1.3.6.1.4.1.3755.0.0.31",
      "name": ".1.3.6.1.2.1.1.2.0",
      "kind": "string",
      "oid": ".1.3.6.1.2.1.1.2.0"
    },
    {
      "value": 460812929,
      "name": ".1.3.6.1.2.1.1.3.0",
      "kind": "int",
      "oid": ".1.3.6.1.2.1.1.3.0"
    },
    {
      "value": "www.atop.com.tw",
      "name": ".1.3.6.1.2.1.1.4.0",
      "kind": "string",
      "oid": ".1.3.6.1.2.1.1.4.0"
    },
    {
      "value": "switch",
      "name": ".1.3.6.1.2.1.1.5.0",
      "kind": "string",
      "oid": ".1.3.6.1.2.1.1.5.0"
    },
    {
      "value": "00-60-E9-27-E3-39",
      "name": "MAC",
      "kind": "string",
      "oid": ""
    },
    {
      "value": "192.168.13.221",
      "name": "IP",
      "kind": "string",
      "oid": ""
    }
  ]
]                                            
```



### TEST gRPC 
Run gRPC service and expose default port(8080)
```shell
$ docker run --rm --network={network-name} -p 8080:8080 {snmpscan-image} run
```

Run snmpsim
```shell
$ docker run -d --rm --network={network-name} {snmpsim-image} 
```
Testing gRPC  

```shell
# Get APIs list
$ grpcurl -plaintext localhost:8080 list scan.SnmpScan
scan.SnmpScan.Describe
scan.SnmpScan.Get
scan.SnmpScan.GetAsyncScanResult
scan.SnmpScan.Scan
scan.SnmpScan.StartAsyncScan
scan.SnmpScan.StopAsyncScan
scan.SnmpScan.Walkall
                                  
# Get 192.168.13.221 snmp information
$ grpcurl -d {\"target\":\"192.168.13.221\"} -plaintext localhost:8080 scan.SnmpScan.Get
{                                             
  "name": ".1.3.6.1.2.1.1.1.0",               
  "value": "Managed Switch, EHG7508-8PoE",    
  "kind": "string"                            
}                                             
{                                             
  "name": ".1.3.6.1.2.1.1.2.0",               
  "value": ".1.3.6.1.4.1.3755.0.0.31",        
  "kind": "string"                            
}                                             
{                                             
  "name": ".1.3.6.1.2.1.1.3.0",               
  "value": "399005846",                       
  "kind": "int"                               
}                                             
{                                             
  "name": ".1.3.6.1.2.1.1.4.0",               
  "value": "www.atop.com.tw",                 
  "kind": "string"                            
}                                             
{                                             
  "name": ".1.3.6.1.2.1.1.5.0",               
  "value": "switch",                          
  "kind": "string"                            
}                                             

# Scan 192.168.13.1/24
$ grpcurl -d {\"range\":\"192.168.13.221/24\"} -plaintext localhost:8080 scan.SnmpScan.Scan                                 
{
  "ip": "192.168.13.200",
  "pdus": [
    {
      "name": ".1.3.6.1.2.1.1.1.0",
      "value": "HP ETHERNET MULTI-ENVIRONMENT,SN:VNBNM37B77,FN:5L5GJTL,SVCID:10288,PID:HP Color LaserJet MFP M281fdw",
      "kind": "string"
    },
    {
      "name": ".1.3.6.1.2.1.1.2.0",
      "value": ".1.3.6.1.4.1.11.2.3.9.1",
      "kind": "string"
    },
    {
      "name": ".1.3.6.1.2.1.1.3.0",
      "value": "96376336",
      "kind": "int"
    },
    {
      "name": ".1.3.6.1.2.1.1.4.0",
      "value": "",
      "kind": "string"
    },
    {
      "name": ".1.3.6.1.2.1.1.5.0",
      "value": "NPI405F13",
      "kind": "string"
    },
    {
      "name": "MAC",
      "value": "C4-65-16-40-5F-13",
      "kind": "string"
    }
  ]
}
{
  "ip": "192.168.13.221",
  "pdus": [
    {
      "name": ".1.3.6.1.2.1.1.1.0",
      "value": "Managed Switch, EHG7508-8PoE",
      "kind": "string"
    },
    {
      "name": ".1.3.6.1.2.1.1.2.0",
      "value": ".1.3.6.1.4.1.3755.0.0.31",
      "kind": "string"
    },
    {
      "name": ".1.3.6.1.2.1.1.3.0",
      "value": "399007760",
      "kind": "int"
    },
    {
      "name": ".1.3.6.1.2.1.1.4.0",
      "value": "www.atop.com.tw",
      "kind": "string"
    },
    {
      "name": ".1.3.6.1.2.1.1.5.0",
      "value": "switch",
      "kind": "string"
    },
    {
      "name": "MAC",
      "value": "00-60-E9-27-E3-39",
      "kind": "string"
    }
  ]
}        

# Scan 192.168.13.1/24 but just for Atop devices
$ grpcurl -d {\"range\":\"192.168.13.221/24\",\"atop_devices\":true} -plaintext localhost:8080 scan.SnmpScan.Scan
{                                                                        
  "ip": "192.168.13.221",                                                
  "pdus": [                                                              
    {                                                                    
      "name": ".1.3.6.1.2.1.1.1.0",                                      
      "value": "Managed Switch, EHG7508-8PoE",                           
      "kind": "string"                                                   
    },                                                                   
    {                                                                    
      "name": ".1.3.6.1.2.1.1.2.0",                                      
      "value": ".1.3.6.1.4.1.3755.0.0.31",                               
      "kind": "string"                                                   
    },                                                                   
    {                                                                    
      "name": ".1.3.6.1.2.1.1.3.0",                                      
      "value": "399016262",                                              
      "kind": "int"                                                      
    },                                                                   
    {                                                                    
      "name": ".1.3.6.1.2.1.1.4.0",                                      
      "value": "www.atop.com.tw",                                        
      "kind": "string"                                                   
    },                                                                   
    {                                                                    
      "name": ".1.3.6.1.2.1.1.5.0",                                      
      "value": "switch",                                                 
      "kind": "string"                                                   
    },                                                                   
    {                                                                    
      "name": "MAC",                                                     
      "value": "00-60-E9-27-E3-39",                                      
      "kind": "string"                                                   
    }                                                                    
  ]                                                                      
}                                                                        


```

### Testing in container  
**Build**  
```shell
$ docker build --no-cache -t {tag} .
```
**Run**  
testing-net is a network which attached to `serviceswatcher` container  

```shell
$ docker run --rm --network testing-net -h snmpscan {tag} {arg}
```

**Async APIs pre-requirements**

- `serviceswatcher` service should running in the same network
- `redis` should running in the same network

```shell
# Start serviceswatcher
$  docker run -d --rm --network testing-net -h serviceswatcher serviceswatcher:latest run
# Start redis
$  docker run -d --rm --network testing-net -h redis redis
```

Async grpc sample

```shell
# scan
$ grpcurl -d {\"range\":\"192.168.13.221/24\",\"atop_devices\":true} -plaintext localhost:8080 scan.SnmpScan.StartAsyncScan
{
  "sessionId": "ss:5b0457b7-cf69-4b58-9f47-8689a729c5e5",
  "success": true
}

# get result
$ ggrpcurl -d {\"session_id\":\"ss:5b0457b7-cf69-4b58-9f47-8689a729c5e5\"} -plaintext localhost:8080 scan.SnmpScan.GetAsyncScanResult
{
  "sessionId": "ss:58e0021f-bc62-49a3-b0fb-76f52f688dfd",
  "success": true,
  "status": "done",
  "result": [
    {
      "ip": "192.168.13.200",
      "pdus": [
        {
          "name": ".1.3.6.1.2.1.1.1.0",
          "value": "HP ETHERNET MULTI-ENVIRONMENT,SN:VNBNM37B77,FN:5L5GJTL,SVCID:10288,PID:HP Color LaserJet MFP M281fdw",
          "kind": "string"
        },
        {
          "name": ".1.3.6.1.2.1.1.2.0",
          "value": ".1.3.6.1.4.1.11.2.3.9.1",
          "kind": "string"
        },
        {
          "name": ".1.3.6.1.2.1.1.3.0",
          "value": 97060861,
          "kind": "int"
        },
        {
          "name": ".1.3.6.1.2.1.1.4.0",
          "value": "",
          "kind": "string"
        },
        {
          "name": ".1.3.6.1.2.1.1.5.0",
          "value": "NPI405F13",
          "kind": "string"
        },
        {
          "name": "MAC",
          "value": "C4-65-16-40-5F-13",
          "kind": "string"
        }
      ]
    },
    {
      "ip": "192.168.13.221",
      "pdus": [
        {
          "name": ".1.3.6.1.2.1.1.1.0",
          "value": "Managed Switch, EHG7508-8PoE",
          "kind": "string"
        },
        {
          "name": ".1.3.6.1.2.1.1.2.0",
          "value": ".1.3.6.1.4.1.3755.0.0.31",
          "kind": "string"
        },
        {
          "name": ".1.3.6.1.2.1.1.3.0",
          "value": 399692294,
          "kind": "int"
        },
        {
          "name": ".1.3.6.1.2.1.1.4.0",
          "value": "www.atop.com.tw",
          "kind": "string"
        },
        {
          "name": ".1.3.6.1.2.1.1.5.0",
          "value": "switch",
          "kind": "string"
        },
        {
          "name": "MAC",
          "value": "00-60-E9-27-E3-39",
          "kind": "string"
        }
      ]
    }
  ]
}
```





### CLI

#### Scan
Scan snmp agents in the classless inter-domain routing 
```shell
# scan 192.168.13.1 - 192.168.13.255
$ snmpscan scan --range=192.168.13.1/24
[
  {
    "ip": "192.168.13.200",
    "data": [
      {
        "value": "HP ETHERNET MULTI-ENVIRONMENT,SN:VNBNM37B77,FN:5L5GJTL,SVCID:10288,PID:HP Color LaserJet MFP M281fdw",
        "name": ".1.3.6.1.2.1.1.1.0",
        "kind": "string"
      },
      {
        "value": ".1.3.6.1.4.1.11.2.3.9.1",
        "name": ".1.3.6.1.2.1.1.2.0",
        "kind": "string"
      },
      {
        "value": 63241880,
        "name": ".1.3.6.1.2.1.1.3.0",
        "kind": "int"
      },
      {
        "value": "",
        "name": ".1.3.6.1.2.1.1.4.0",
        "kind": "string"
      },
      {
        "value": "NPI405F13",
        "name": ".1.3.6.1.2.1.1.5.0",
        "kind": "string"
      },
      {
        "value": "C4-65-16-40-5F-13",
        "name": "MAC",
        "kind": "string"
      }
    ]
  },
  {
    "ip": "192.168.13.221",
    "data": [
      {
        "value": "Managed Switch, EHG7508-8PoE",
        "name": ".1.3.6.1.2.1.1.1.0",
        "kind": "string"
      },
      {
        "value": ".1.3.6.1.4.1.3755.0.0.31",
        "name": ".1.3.6.1.2.1.1.2.0",
        "kind": "string"
      },
      {
        "value": 365872858,
        "name": ".1.3.6.1.2.1.1.3.0",
        "kind": "int"
      },
      {
        "value": "www.atop.com.tw",
        "name": ".1.3.6.1.2.1.1.4.0",
        "kind": "string"
      },
      {
        "value": "switch",
        "name": ".1.3.6.1.2.1.1.5.0",
        "kind": "string"
      },
      {
        "value": "00-60-E9-27-E3-39",
        "name": "MAC",
        "kind": "string"
      }
    ]
  }
]
```
Or you can just scan Atop devices in the CIDR with `--atop-devices` flag

```shell
$ snmpscan scan --range 192.168.13.1/24 --atop-devices                                  
E0418 20:07:57.398944 31792 utils.go:34] OIDs file not found, use default oids instead of      
[                                                                                              
  {                                                                                            
    "ip": "192.168.13.221",                                                                    
    "data": [                                                                                  
      {                                                                                        
        "value": "Managed Switch, EHG7508-8PoE",                                               
        "name": ".1.3.6.1.2.1.1.1.0",                                                          
        "kind": "string"                                                                       
      },                                                                                       
      {                                                                                        
        "value": ".1.3.6.1.4.1.3755.0.0.31",                                                   
        "name": ".1.3.6.1.2.1.1.2.0",                                                          
        "kind": "string"                                                                       
      },                                                                                       
      {                                                                                        
        "value": 365882500,                                                                    
        "name": ".1.3.6.1.2.1.1.3.0",                                                          
        "kind": "int"                                                                          
      },                                                                                       
      {                                                                                        
        "value": "www.atop.com.tw",                                                            
        "name": ".1.3.6.1.2.1.1.4.0",                                                          
        "kind": "string"                                                                       
      },                                                                                       
      {                                                                                        
        "value": "switch",                                                                     
        "name": ".1.3.6.1.2.1.1.5.0",                                                          
        "kind": "string"                                                                       
      },                                                                                       
      {                                                                                        
        "value": "00-60-E9-27-E3-39",                                                          
        "name": "MAC",                                                                         
        "kind": "string"                                                                       
      }                                                                                        
    ]                                                                                          
  }                                                                                            
]                                                                                              
```

#### Describe
Describe particular device

```shell
$ snmpscan describe 192.168.13.221
{
  "sysDescr": "Managed Switch, EHG7508-8PoE",
  "sysUpTime": 365899159,
  "sysContact": "www.atop.com.tw",
  "sysName": "switch",
  "sysLocation": "Switch's Location",
  "interfaces": [
    {
      "index": "1",
      "description": "lo",
      "MAC": "00-60-E9-27-E3-39",
      "portEnable": 1,
      "portStatus": 1,
      "inNUcastPkts": 0,
      "inErrors": 0,
      "outNUcastPkts": 0,
      "outErrors": 0,
      "inOctets": 0,
      "InUcastPkts": 0,
      "inMulticastPkts": 0,
      "inBroadcastPkts": 0,
      "outOctets": 0,
      "outUcastPkts": 0,
      "outMulticastPkts": 0,
      "outBroadcastPkts": 0,
      "highSpeed": 10
    },
    {
      "index": "2",
      "description": "eth1",
      "MAC": "00-60-E9-27-E3-39",
      "portEnable": 1,
      "portStatus": 1,
      "inNUcastPkts": 0,
      "inErrors": 0,
      "outNUcastPkts": 0,
      "outErrors": 0,
      "inOctets": 0,
      "InUcastPkts": 0,
      "inMulticastPkts": 0,
      "inBroadcastPkts": 0,
      "outOctets": 0,
      "outUcastPkts": 0,
      "outMulticastPkts": 0,
      "outBroadcastPkts": 0,
      "highSpeed": 10
    },
    {
      "index": "3",
      "description": "Port3",
      "MAC": "00-60-E9-27-E3-39",
      "portEnable": 1,
      "portStatus": 2,
      "inNUcastPkts": 0,
      "inErrors": 0,
      "outNUcastPkts": 0,
      "outErrors": 0,
      "inOctets": 0,
      "InUcastPkts": 0,
      "inMulticastPkts": 0,
      "inBroadcastPkts": 0,
      "outOctets": 0,
      "outUcastPkts": 0,
      "outMulticastPkts": 0,
      "outBroadcastPkts": 0,
      "highSpeed": 0
    },
    {
      "index": "4",
      "description": "Port4",
      "MAC": "00-60-E9-27-E3-39",
      "portEnable": 1,
      "portStatus": 1,
      "inNUcastPkts": 4069165,
      "inErrors": 0,
      "outNUcastPkts": 4232201,
      "outErrors": 0,
      "inOctets": 794696189,
      "InUcastPkts": 7328610,
      "inMulticastPkts": 239,
      "inBroadcastPkts": 915910,
      "outOctets": 2080623640,
      "outUcastPkts": 2939889,
      "outMulticastPkts": 3153016,
      "outBroadcastPkts": 1079185,
      "highSpeed": 1000
    },
    {
      "index": "5",
      "description": "Port5",
      "MAC": "00-60-E9-27-E3-39",
      "portEnable": 1,
      "portStatus": 2,
      "inNUcastPkts": 0,
      "inErrors": 0,
      "outNUcastPkts": 0,
      "outErrors": 0,
      "inOctets": 0,
      "InUcastPkts": 0,
      "inMulticastPkts": 0,
      "inBroadcastPkts": 0,
      "outOctets": 0,
      "outUcastPkts": 0,
      "outMulticastPkts": 0,
      "outBroadcastPkts": 0,
      "highSpeed": 0
    },
    {
      "index": "6",
      "description": "Port6",
      "MAC": "00-60-E9-27-E3-39",
      "portEnable": 1,
      "portStatus": 2,
      "inNUcastPkts": 0,
      "inErrors": 0,
      "outNUcastPkts": 0,
      "outErrors": 0,
      "inOctets": 0,
      "InUcastPkts": 0,
      "inMulticastPkts": 0,
      "inBroadcastPkts": 0,
      "outOctets": 0,
      "outUcastPkts": 0,
      "outMulticastPkts": 0,
      "outBroadcastPkts": 0,
      "highSpeed": 0
    },
    {
      "index": "7",
      "description": "Port7",
      "MAC": "00-60-E9-27-E3-39",
      "portEnable": 1,
      "portStatus": 2,
      "inNUcastPkts": 38898,
      "inErrors": 0,
      "outNUcastPkts": 83454,
      "outErrors": 0,
      "inOctets": 606151,
      "InUcastPkts": 2307,
      "inMulticastPkts": 1726,
      "inBroadcastPkts": 717,
      "outOctets": 21966741,
      "outUcastPkts": 781,
      "outMulticastPkts": 36455,
      "outBroadcastPkts": 46999,
      "highSpeed": 0
    },
    {
      "index": "8",
      "description": "Port8",
      "MAC": "00-60-E9-27-E3-39",
      "portEnable": 1,
      "portStatus": 2,
      "inNUcastPkts": 0,
      "inErrors": 0,
      "outNUcastPkts": 0,
      "outErrors": 0,
      "inOctets": 0,
      "InUcastPkts": 0,
      "inMulticastPkts": 0,
      "inBroadcastPkts": 0,
      "outOctets": 0,
      "outUcastPkts": 0,
      "outMulticastPkts": 0,
      "outBroadcastPkts": 0,
      "highSpeed": 0
    },
    {
      "index": "10000",
      "description": "sit0",
      "MAC": "00-00-00-00-E3-39",
      "portEnable": 2,
      "portStatus": 1,
      "inNUcastPkts": 0,
      "inErrors": 0,
      "outNUcastPkts": 0,
      "outErrors": 0,
      "inOctets": 0,
      "InUcastPkts": 0,
      "inMulticastPkts": 0,
      "inBroadcastPkts": 0,
      "outOctets": 0,
      "outUcastPkts": 0,
      "outMulticastPkts": 0,
      "outBroadcastPkts": 0,
      "highSpeed": 0
    },
    {
      "index": "10001",
      "description": "lo",
      "MAC": "",
      "portEnable": 1,
      "portStatus": 1,
      "inNUcastPkts": 0,
      "inErrors": 0,
      "outNUcastPkts": 0,
      "outErrors": 0,
      "inOctets": 0,
      "InUcastPkts": 0,
      "inMulticastPkts": 0,
      "inBroadcastPkts": 0,
      "outOctets": 0,
      "outUcastPkts": 0,
      "outMulticastPkts": 0,
      "outBroadcastPkts": 0,
      "highSpeed": 10
    },
    {
      "index": "10002",
      "description": "eth1",
      "MAC": "00-60-E9-27-E3-39",
      "portEnable": 1,
      "portStatus": 1,
      "inNUcastPkts": 0,
      "inErrors": 0,
      "outNUcastPkts": 0,
      "outErrors": 0,
      "inOctets": 0,
      "InUcastPkts": 0,
      "inMulticastPkts": 0,
      "inBroadcastPkts": 0,
      "outOctets": 0,
      "outUcastPkts": 0,
      "outMulticastPkts": 0,
      "outBroadcastPkts": 0,
      "highSpeed": 10
    }
  ]
}
```