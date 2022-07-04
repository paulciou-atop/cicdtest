## logging





## Examples

#### server 

How import the package

```go
import AtopSyslog "atopsyslog/syslog"
```

 File name: syslog.log

```go
import AtopSyslog "atopsyslog/syslog"
func main() {
	server := AtopSyslog.NewServer("syslog.log")
	err := server.Run("0.0.0.0:514")
	if err != nil {
		log.Fatal(err)
	}
}

```



syslog.log record

```json
{"client":"127.0.0.1:63937","content":"test message","facility":5,"hostname":"TE-YanLin-NB10","priority":47,"severity":7,"tag":"demotag","timestamp":"2022-03-18T10:32:19+08:00","tls_peer":""}
```



#### client

```go
sysLog, err := AtopSyslog.DialLogger("udp", "0.0.0.0:514", AtopSyslog.LOG_DEBUG|AtopSyslog.LOG_SYSLOG, "demotag")
	if err != nil {
		log.Fatal(err)
	}
	//	fmt.Fprintf(sysLog, "This is a daemon warning with demotag.")
	sysLog.Write([]byte("test message"))
	sysLog.Close()
```

