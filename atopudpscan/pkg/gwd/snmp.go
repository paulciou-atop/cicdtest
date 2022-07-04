package gwd

import (
	"bytes"
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

type snmpItem = []byte
type Snmp struct {
	Name     string
	Location string
	Contact  string
}

const (
	iac             = 0xff
	sb              = 0xfa
	com_port_option = 0x2c
	signature_host  = 0x00
	se              = 0xf0
	//SNMP_HOST       = 89
)

var (
	QuerySnmp_name     snmpItem = []byte{0x2c, 0x59, 0x2}
	QuerySnmp_location snmpItem = []byte{0x2c, 0x59, 0x4}
	QuerySnmp_contact  snmpItem = []byte{0x2c, 0x59, 0x6}

	ParsingSnmp_name     snmpItem = []byte{0x2c, 0xbd, 0x3}
	ParsingSnmp_location snmpItem = []byte{0x2c, 0xbd, 0x5}
	ParsingSnmp_contact  snmpItem = []byte{0x2c, 0xbd, 0x7}

	SettingSnmp_name    snmpItem = []byte{0x2c, 0x59, 0x3}
	SettinSnmp_location snmpItem = []byte{0x2c, 0x59, 0x5}
	SettinSnmp_contact  snmpItem = []byte{0x2c, 0x59, 0x7}
)

func (d *Device) QuerySnmp(password string) (*Snmp, error) {
	s := &Snmp{}
	packet := QuerynmpPacket(password)
	address := strings.Join([]string{d.ip, strconv.Itoa(loginport)}, ":")
	con, err := net.DialTimeout("tcp", address, conntimeout)
	if err != nil {
		return nil, err
	}
	defer con.Close()
	dst := make([]byte, 1024)
	con.SetReadDeadline(time.Now().Add(3 * time.Second))
	con.Read(dst)
	_, err = con.Write(packet)
	if err != nil {

		return nil, err
	}
	con.SetReadDeadline(time.Now().Add(3 * time.Second))
	dst = make([]byte, 1024)
	n, err := con.Read(dst)

	if err != nil {

		return nil, err
	}
	dst = dst[:n]
	r, err := parseinglogin(dst, Login)
	if err != nil {

		return nil, err
	}
	if !r {
		return nil, errors.New("password error")
	}

	s.Name, err = parseingSnmp(dst, ParsingSnmp_name)
	if err != nil {

		return nil, err
	}
	s.Contact, err = parseingSnmp(dst, ParsingSnmp_contact)
	if err != nil {

		return nil, err
	}
	s.Location, err = parseingSnmp(dst, ParsingSnmp_location)
	if err != nil {

		return nil, err
	}

	return s, nil

}

func (d *Device) SettingSnmp(password string, s *Snmp) (*Snmp, error) {
	packet := settingSnmpPacket(password, s)
	address := strings.Join([]string{d.ip, strconv.Itoa(loginport)}, ":")
	con, err := net.DialTimeout("tcp", address, conntimeout)
	if err != nil {
		return nil, err
	}
	defer con.Close()
	dst := make([]byte, 1024)
	con.SetReadDeadline(time.Now().Add(3 * time.Second))
	con.Read(dst)
	_, err = con.Write(packet)
	if err != nil {

		return nil, err
	}
	con.SetReadDeadline(time.Now().Add(3 * time.Second))
	dst = make([]byte, 1024)
	n, err := con.Read(dst)

	if err != nil {

		return nil, err
	}
	dst = dst[:n]
	r, err := parseinglogin(dst, Login)
	if err != nil {

		return nil, err
	}
	if !r {
		return nil, errors.New("password error")
	}

	s.Name, err = parseingSnmp(dst, ParsingSnmp_name)
	if err != nil {

		return nil, err
	}
	s.Contact, err = parseingSnmp(dst, ParsingSnmp_contact)
	if err != nil {

		return nil, err
	}
	s.Location, err = parseingSnmp(dst, ParsingSnmp_location)
	if err != nil {

		return nil, err
	}

	return s, nil
}
func settingSnmpPacket(password string, s *Snmp) []byte {
	packet := loginPacket(password)

	packet = append(packet, settingSnmpNamePacket(s.Name)...)

	packet = append(packet, settingSnmpLocationPacket(s.Location)...)

	packet = append(packet, settingSnmpContactPacket(s.Contact)...)

	return packet
}

func settingSnmpNamePacket(s string) []byte {
	packet := make([]byte, 0)
	packet = append(packet, iac)
	packet = append(packet, sb)
	packet = append(packet, SettingSnmp_name...)
	packet = append(packet, s...)
	packet = append(packet, 0x00)
	packet = append(packet, iac)
	packet = append(packet, se)
	return packet
}
func settingSnmpLocationPacket(s string) []byte {
	packet := make([]byte, 0)
	packet = append(packet, iac)
	packet = append(packet, sb)
	packet = append(packet, SettinSnmp_location...)
	packet = append(packet, s...)
	packet = append(packet, 0x00)
	packet = append(packet, iac)
	packet = append(packet, se)
	return packet
}

func settingSnmpContactPacket(s string) []byte {
	packet := make([]byte, 0)
	packet = append(packet, iac)
	packet = append(packet, sb)
	packet = append(packet, SettinSnmp_contact...)
	packet = append(packet, s...)
	packet = append(packet, 0x00)
	packet = append(packet, iac)
	packet = append(packet, se)
	return packet
}

func QuerynmpPacket(password string) []byte {
	packet := loginPacket(password)

	packet = append(packet, querysnmpNamePacket()...)

	packet = append(packet, querysnmpLocationPacket()...)

	packet = append(packet, querysnmpContactPacket()...)

	return packet
}
func querysnmpNamePacket() []byte {
	packet := make([]byte, 0)
	packet = append(packet, iac)
	packet = append(packet, sb)
	packet = append(packet, QuerySnmp_name...)
	packet = append(packet, iac)
	packet = append(packet, se)
	return packet
}

func querysnmpLocationPacket() []byte {
	packet := make([]byte, 0)
	packet = append(packet, iac)
	packet = append(packet, sb)
	packet = append(packet, QuerySnmp_location...)
	packet = append(packet, iac)
	packet = append(packet, se)
	return packet
}

func querysnmpContactPacket() []byte {
	packet := make([]byte, 0)
	packet = append(packet, iac)
	packet = append(packet, sb)
	packet = append(packet, QuerySnmp_contact...)
	packet = append(packet, iac)
	packet = append(packet, se)
	return packet
}

func parseingSnmp(b []byte, value snmpItem) (string, error) {
	for {
		index := bytes.Index(b, []byte{iac, sb})
		last := bytes.Index(b, []byte{iac, se})
		v := b[index+2 : last]
		index = bytes.Index(v, value)
		if index >= 0 {
			return string(v[len(value):]), nil
		}
		b = RemoveIndex(b, last+2)
		if len(b) <= 0 {
			return "", errors.New("parsing error,value not exist")
		}
	}

}

func RemoveIndex(s []byte, index int) []byte {
	if len(s)-index >= 0 {
		return s[index:]
	}
	return s
}
