package modules

// DeviceResult -> device_results
type DeviceResult struct {
	ID          int64
	SessionID   string
	Model       string
	MacAddress  string
	IpAddress   string
	Netmask     string
	Gateway     string
	Hostname    string
	Kernel      string
	Ap          string
	FirmwareVer string
	Description string
}

// DeviceSession -> device_sessions
type DeviceSession struct {
	ID              int64  //id
	SessionID       string // session_id
	State           string // state
	CreatedTime     string // created_time
	LastUpdatedTime string
}
