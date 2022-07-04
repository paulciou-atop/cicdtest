package gwd

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ModelInfo struct {
	Model      string `json:"model"`
	MACAddress string `json:"macAddress"`
	IPAddress  string `json:"iPAddress"`
	Netmask    string `json:"netmask"`
	Gateway    string `json:"gateway"`
	Hostname   string `json:"hostname"`
	Kernel     string `json:"kernel"`
	Ap         string `json:"ap"`
	IsDHCP     bool   `json:"isDHCP"`
}

type NetworkConfig struct {
	MACAddress   string `json:"mACAddress"`
	IPAddress    string `json:"iPAddress"`
	NewIPAddress string `json:"newIPAddress"`
	Netmask      string `json:"netmask"`
	Gateway      string `json:"gateway"`
	Hostname     string `json:"hostname"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

const Port = 55954

//localnetwork:out path  form  local ip
func NewAtopGwd(localnetwork string) *AtopGwd {
	return &AtopGwd{network: localnetwork}
}

type AtopGwd struct {
	network string
}

func (a *AtopGwd) Scanner() error {
	return a.broadCast(invitePacket())
}

func (a *AtopGwd) Beep(n NetworkConfig) error {
	b, err := settingBeepPacket(n)
	if err != nil {
		return err
	}
	return a.broadCast(b)
}

func (a *AtopGwd) Reboot(n NetworkConfig) error {
	b, err := settingRebootPacket(n)
	if err != nil {
		return err
	}
	return a.broadCast(b)
}

func (a *AtopGwd) ResetToDefault(n NetworkConfig) error {
	b, err := settingResetDefault(n)
	if err != nil {
		return err
	}
	return a.broadCast(b)
}

func (a *AtopGwd) SettingConfig(n NetworkConfig) error {
	b, err := settingConfigPacket(n)
	if err != nil {
		return err
	}
	return a.broadCast(b)
}

//ConfigPacket
func settingConfigPacket(config NetworkConfig) ([]byte, error) {
	packet := configPacket()
	packet, err := addIpToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addNewIpToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addNewGatewayToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addNewNetmaskToPackage(config, packet)
	if err != nil {
		return nil, err
	}

	packet, err = addMacToPackage(config, packet)
	if err != nil {
		return nil, err
	}

	packet, err = addhostNameToPackage(config, packet)
	if err != nil {
		return nil, err
	}

	packet = addUserAndPwdToPackage(config, packet)
	return packet, nil
}

//RebootPacket
func settingRebootPacket(config NetworkConfig) ([]byte, error) {
	packet := rebootPacket()
	packet, err := addIpToPackage(config, packet)
	if err != nil {
		return nil, errors.New("IPAddress format error")
	}
	packet, err = addMacToPackage(config, packet)
	if err != nil {
		return nil, fmt.Errorf("MacAddress format error ,reason:%s", err.Error())
	}
	packet = addUserAndPwdToPackage(config, packet)
	return packet, nil
}

//BeepPacket
func settingBeepPacket(config NetworkConfig) ([]byte, error) {
	packet := beepPacket()
	packet, err := addIpToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addMacToPackage(config, packet)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

func settingResetDefault(config NetworkConfig) ([]byte, error) {
	packet := reSetDefaultPacket()
	packet, err := addIpToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addMacToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet = addUserAndPwdToPackage(config, packet)
	return packet, nil
}

//BroadCast
func (n *AtopGwd) broadCast(msg []byte) error {
	address := strings.Join([]string{"255.255.255.255", strconv.Itoa(Port)}, ":")
	broadcastAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	local, err := net.ResolveUDPAddr("udp", n.network)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", local, broadcastAddr)
	if err != nil {
		return err
	}
	_, err = conn.Write(msg)
	defer conn.Close()
	return err
}

//addIpToPackage
func addIpToPackage(config NetworkConfig, packet []byte) ([]byte, error) {
	ip, err := parsingIp(config.IPAddress)
	if err != nil {
		return nil, errors.New("IPAddress format error")
	}
	for i := 0; i < 4; i++ {
		packet[12+i] = ip[i]
	}
	return packet, nil
}

//add ip changed to package
func addNewIpToPackage(config NetworkConfig, packet []byte) ([]byte, error) {
	ip, err := parsingIp(config.NewIPAddress)
	if err != nil {
		return nil, errors.New("new IPAddress format error")
	}
	for i := 0; i < 4; i++ {
		packet[16+i] = ip[i]
	}
	return packet, nil
}

//add Netmask to package
func addNewNetmaskToPackage(config NetworkConfig, packet []byte) ([]byte, error) {
	ip, err := parsingIp(config.Netmask)
	if err != nil {
		return nil, errors.New("netmask format error")
	}
	for i := 0; i < 4; i++ {
		packet[236+i] = ip[i]
	}
	return packet, nil
}

//add Gateway to package
func addNewGatewayToPackage(config NetworkConfig, packet []byte) ([]byte, error) {
	ip, err := parsingIp(config.Gateway)
	if err != nil {
		return nil, errors.New("gateway format error")
	}
	for i := 0; i < 4; i++ {
		packet[24+i] = ip[i]
	}
	return packet, nil
}

//add host Name To Package
func addhostNameToPackage(config NetworkConfig, packet []byte) ([]byte, error) {
	i := 90
	host := []byte(config.Hostname)
	h, err := EncodeBig5(host)
	if err == nil {
		host = h
	}
	if len(host) >= 16 {
		return nil, fmt.Errorf("HostName len is too long")
	}
	for _, v := range host {
		packet[i] = v
		i++
	}
	return packet, nil
}

//add Mac To Package
func addMacToPackage(config NetworkConfig, packet []byte) ([]byte, error) {
	mac, err := parsingMacAddress(config.MACAddress)
	if err != nil {
		return nil, fmt.Errorf("MacAddress format error ,reason:%s", err.Error())
	}
	for i := 0; i < 6; i++ {
		packet[28+i] = mac[i]
	}
	return packet, nil
}

//addUserAndPwdtoPackage
func addUserAndPwdToPackage(config NetworkConfig, packet []byte) []byte {
	packet[70] = 2
	i := 71
	user := []byte(config.Username)
	pwd := []byte(config.Password)
	for _, v := range user {
		packet[i] = v
		i++
	}
	packet[i] = 0x20
	i++
	for _, v := range pwd {
		packet[i] = v
		i++
	}
	return packet
}
