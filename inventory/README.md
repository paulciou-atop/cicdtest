Inventory types and packages used in the NMS.

## device
Device inventory, this service try to gather all scan result from other scanner services. Polished, render and merge data into database.

To test, build this service, the service needs 
- message queue service, the service subscribe scan results to get last device information.
- data store, e.g. postgresql, mongo.  the service write polished data into database. 

### Testing

#### List inventories
```shell
$ grpcurl -d '{"pagination":{"page":1, "size":20}}' -plaintext localhost:8101 inventory.Inventories.List
{
  "success": true,
  "pagination": {
    "page": 1,
    "size": 20,
    "total": 1
  },
  "inventories": [
    {
      "deviceType": "snmp",
      "id": "60f2d4cc-6f22-4702-bacc-1cd26bade276",
      "model": "EHG7508-8PoE",
      "location": {

      },
      "ipAddress": "192.168.13.221",
      "macAddress": "00-60-E9-27-E3-39",
      "firmwareInformation": {
        "kernel": "70.09"
      },
      "createdAt": "2022-06-22 09:51:39",
      "lastSeen": "2022-06-22 09:51:53",
      "supportProtocols": [
        "snmp"
      ],
      "more": {
        }
    }
  ]
}
```
**HTTP**
```shell
$ curl -d '{"pagination":{"page":1,"size":20}}' -X POST http://localhost:8111/v1/inventories | json_pp -json_opt pretty,canonical
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   544  100   509  100    35   248k  17500 --:--:-- --:--:-- --:--:--  265k
{
   "inventories" : [
      {
         "createdAt" : "2022-06-22 09:51:39",
         "deviceType" : "snmp",
         "firmwareInformation" : {
            "ap" : "",
            "kernel" : "70.09"
         },
         "hostName" : "",
         "id" : "60f2d4cc-6f22-4702-bacc-1cd26bade276",
         "ipAddress" : "192.168.13.221",
         "lastMissing" : "",
         "lastRecovered" : "",
         "lastSeen" : "2022-06-22 09:51:53",
         "location" : {
            "path" : ""
         },
         "macAddress" : "00-60-E9-27-E3-39",
         "model" : "EHG7508-8PoE",
         "more" : {},
         "name" : "",
         "owner" : "",
         "supportProtocols" : [
            "snmp"
         ]
      }
   ],
   "message" : "",
   "pagination" : {
      "page" : 1,
      "size" : 20,
      "total" : 1
   },
   "success" : true
}

```

#### Get inventory

**gRPC**

```shell
$ grpcurl -d '{"id":"60f2d4cc-6f22-4702-bacc-1cd26bade276"}' -plaintext localhost:8101 inventory.Inventories.Get
{
  "success": true,
  "inventory": {
    "deviceType": "snmp",
    "id": "60f2d4cc-6f22-4702-bacc-1cd26bade276",
    "model": "EHG7508-8PoE",
    "location": {

    },
    "ipAddress": "192.168.13.221",
    "macAddress": "00-60-E9-27-E3-39",
    "firmwareInformation": {
      "kernel": "70.09"
    },
    "createdAt": "2022-06-22 09:51:39",
    "lastSeen": "2022-06-22 09:51:53",
    "supportProtocols": [
      "snmp"
    ],
    "more": {
      }
  }
}

```