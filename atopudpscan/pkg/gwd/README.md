## ATOP UDP Scan service

Send out UDP broadcast to find ATOP devices on the directly attached subnet



## Examples

How import the package

```go
import 	udp "atopudpscan/Udp"
```

## Gwd

How import the package

```go
import 	"atopudpscan/gwd"
```

### Api List

```go
func NewAtopGwd(localnetwork string) *atopGwd 
func (a *atopGwd) Scanner() error 
func (a *atopGwd) Beep(n NetworkConfig) error 
func (a *atopGwd) Reboot(n NetworkConfig) error 
func (a *atopGwd) ResetToDefault(n NetworkConfig) error
func (a *atopGwd) SettingConfig(n NetworkConfig) error
```

#### NewAtopGwd

```go
	g := gwd.NewAtopGwd("192.168.4.21:0")
```

#### Scanner

```go
func TestScanner(t *testing.T) {

	g := gwd.NewAtopGwd("192.168.4.21:0")
	err := g.Scanner()
	if err != nil {

		t.Error(err)
	}
}
```

#### Beep

```go

func TestBeep(t *testing.T) {
	net := gwd.NetworkConfig{
		IPAddress:  "192.168.4.30",
		MACAddress: "00:60:e9:11:22:33",
	}

	g := gwd.NewAtopGwd("192.168.4.21:0")
	err := g.Beep(net)
	if err != nil {

		t.Error(err)
	}
}
```



#### Reboot

```go
func TestReboot(t *testing.T) {
	net := gwd.NetworkConfig{
		IPAddress:  "192.168.4.25",
		MACAddress: "00:60:e9:11:22:33",
		Username:   "admin",
		Password:   "default",
	}

	g := gwd.NewAtopGwd("192.168.4.21:0")
	err := g.Reboot(net)
	if err != nil {

		t.Error(err)
	}
}
```

#### ResetToDefault

```go
func TestResetToDefault(t *testing.T) {
	net := gwd.NetworkConfig{
		IPAddress:  "192.168.4.30",
		MACAddress: "00:60:e9:11:22:33",
		Username:   "admin",
		Password:   "default",
	}

	g := gwd.NewAtopGwd("192.168.4.21:0")
	err := g.ResetToDefault(net)
	if err != nil {

		t.Error(err)
	}
}
```

####  SettingConfig

```go
func TestSettingConfig(t *testing.T) {
	net := gwd.NetworkConfig{
		IPAddress:    "192.168.4.30",
		MACAddress:   "00:60:e9:11:22:33",
		NewIPAddress: "192.168.4.30",
		Netmask:      "255.255.255.0",
		Gateway:      "192.168.4.254",
		Hostname:     "atoptest",
		Username:     "admin",
		Password:     "default",
	}

	g := gwd.NewAtopGwd("192.168.4.21:0")
	err := g.SettingConfig(net)
	if err != nil {

		t.Error(err)
	}
}

```



## Control Device

How import the package

```go
import 	"atopudpscan/gwd"
```

### Api List

```go
func NewDevice(ip string) *device 
```

### Firmware

#### Api List

```go
func (d *device) FirmWare() *firmWare 
func (f *firmWare) Upgrading(file *os.File) error 
func (f *firmWare) GetProcessStatus() string //return status
```

#####  FirmWare()

```go
device := gwd.NewDevice("192.168.4.30").FirmWare()
```

##### Upgrading(),

##### GetProcessStatus()

```go
func TestFwUpgrading(t *testing.T) {
	f, err := os.Open("SE52XX_K108A120.dld")
	if err != nil {

		t.Error(err)
		return
	}
	defer f.Close()

	device := gwd.NewDevice("192.168.4.30").FirmWare()
	go func() {
		for {
			time.Sleep(300)
			s := device.GetProcessStatus()
			fmt.Print(s)
			if s == gwd.Complete {
				break
			}

		}
	}()
	err = device.Upgrading(f)
	if err != nil {

		t.Error(err)
	}

}
```

###### ProcessStatus 

```go
	process //if in file transfer,retrun 0~100
	None      
	Upgrading 
	Complete  
```

### Login

#### Api List

```go
func (d *device) Login(password string) (bool, error) 
```

##### Login()

```go
func TestLogin(t *testing.T) {
	device := gwd.NewDevice("192.168.4.30")
	r, err := device.Login("default")
	if err != nil {

		t.Error(err)
	}
	if r {
		t.Log("pass")

	} else {
		t.Error("faield")

	}

}
```

### Snmp

#### Api List

```go
func (d *device) QuerySnmp(password string) (*Snmp, error)
func (d *device) SettingSnmp(password string, s *Snmp) (*Snmp, error) 
```

##### QuerySnmp()

```go
	device := gwd.NewDevice("192.168.4.30")

	s, err := device.QuerySnmp("default")
	if err != nil {

		t.Error(err)
		return
	}
	fmt.Print(s.Name)
	fmt.Print(s.Location)
	fmt.Print(s.Contact)
```

##### SettingSnmp()

```go
func TestSettingSnmp(t *testing.T) {
	device := gwd.NewDevice("192.168.4.30")
	snmp := &gwd.Snmp{Contact: "test", Name: "52t789", Location: "taivh"}

	s, err := device.SettingSnmp("default", snmp)
	if err != nil {

		t.Error(err)
		return
	}
	fmt.Print(s.Name)
	fmt.Print(s.Location)
	fmt.Print(s.Contact)
}
```

