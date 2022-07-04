## ATOP UDP Scan service

Send out UDP broadcast to find ATOP devices on the directly attached subnet



## Examples

How import the package

```go
import 	udp "atopudpscan/Udp"
```

## Server

### Api List

```go
func NewAtopUdpServer(ip string) *AtopUdpServer //Get Udp server
func (r *AtopUdpServer) Run() error //run server
func (r *AtopUdpServer) GetReceiveData() ([]byte, error) //get receive data from array,every get ,it clear array
func (r *AtopUdpServer) Stop() //server stop
func (r *AtopUdpServer) SaveToDatabase(b []byte) error 
func (r *AtopUdpServer) GetDataFromDatabase() ([]byte, error) 
func (r *AtopUdpServer) CleanDatabase() error 
```

#### Run()

```go
func main() {
	s := udp.NewAtopUdpServer("0.0.0.0")
	go s.Run()
	for {
		v, err := s.GetReceiveData()
		if err == nil {
			fmt.Println(string(v))
		}
	}

}
```

#### GetReceiveData()

```GO
		v, err := s.GetReceiveData()
		if err == nil {
			fmt.Println(string(v))
		}
```

###### Response:

```json
[{"model":"SE9501","macAddress":"00-60-E9-04-23-01","iPAddress":"192.168.4.193","netmask":"255.255.0","gateway":"192.168.4.254","hostname":"MW826STH","kernel":"1.0","ap":"MW826STH V0.14 CFG-V0.03","isDHCP":false,"time":"2022-03-31 13:14:34"}]
[{"model":"IO5202","macAddress":"00-60-E9-2C-9B-1D","iPAddress":"192.168.4.71","netmask":"255.255.255","gateway":"192.168.4.254","hostname":"System","kernel":"1.22","ap":"IO5202 V1.25","isDHCP":false,"time":"2022-03-31 13:14:34"}]
```

####  Stop()

```go
s := udp.NewAtopUdpServer("0.0.0.0")
s.Stop()
```

#### SaveToDatabase()

```go
func TestSaveToDatabase(t *testing.T) {
	s := NewAtopUdpServer("0.0.0.0")
	go s.Run()
	go func() {
		for {
			v, err := s.GetReceiveData()
			if err == nil {
				log.Println(string(v))
				err := s.SaveToDatabase(v)
				if err != nil {
					t.Error(err)
					return
				}
			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	s.Stop()
}
```

#### GetDataFromDatabase()

```go
func TestGetDataFromDatabase(t *testing.T) {
	s := NewAtopUdpServer("0.0.0.0")
	b, err := s.GetDataFromDatabase()
	if err == nil {
		log.Println(string(b))
	} else {
		t.Error(err)
	}
}
```

#### CleanDatabase()

```go

func TestCleanDatabase(t *testing.T) {
	s := NewAtopUdpServer("0.0.0.0")
	err := s.CleanDatabase()
	if err != nil {
		t.Error(err)
	}
}

```


