package gwd

import (
	"bytes"
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	loginport = 55960
)
const conntimeout = 10 * time.Second

type LoginItem = []byte

var (
	Login LoginItem = []byte{0x2c, 0x64}
)

func loginPacket(password string) []byte {
	packet := make([]byte, 0)
	packet = append(packet, iac)
	packet = append(packet, sb)
	packet = append(packet, com_port_option)
	packet = append(packet, signature_host)
	for _, v := range password {
		packet = append(packet, byte(v))
	}
	packet = append(packet, iac)
	packet = append(packet, se)
	return packet
}

func (d *Device) Login(password string) (bool, error) {
	packet := loginPacket(password)
	address := strings.Join([]string{d.ip, strconv.Itoa(loginport)}, ":")
	con, err := net.DialTimeout("tcp", address, conntimeout)
	if err != nil {
		return false, err
	}
	defer con.Close()
	dst := make([]byte, 7)
	con.SetReadDeadline(time.Now().Add(3 * time.Second))
	con.Read(dst)
	_, err = con.Write(packet)
	if err != nil {

		return false, err
	}

	con.SetReadDeadline(time.Now().Add(3 * time.Second))
	dst = make([]byte, 1024)
	n, err := con.Read(dst)

	if err != nil {

		return false, err
	}
	dst = dst[:n]
	return parseinglogin(dst, Login)

}

func parseinglogin(b []byte, value LoginItem) (bool, error) {
	for {
		index := bytes.Index(b, []byte{iac, sb})
		last := bytes.Index(b, []byte{iac, se})
		v := b[index+2 : last]
		index = bytes.Index(v, value)
		if index >= 0 {
			if string(v[index+2:]) == "S" {
				return true, nil
			} else {
				return false, nil
			}
		}
		b = RemoveIndex(b, last+2)
		if len(b) <= 0 {
			return false, errors.New("parsing error,value not exist")
		}
	}

}
