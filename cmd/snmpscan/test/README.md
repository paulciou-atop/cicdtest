## SNMP Command Testing


### Pre requirement
Go version : >1.18  

### Start Dummy SNMP server and test

```shell
$ cd test/server
$ go run snmpserver.go
$ cd ..
$ go test
```

