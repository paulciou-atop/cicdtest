package Simulatedevice

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"nms/atopudpscan/pkg/gwd"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

func byteToHexString(msg []byte, sep string) string {
	str := make([]string, len(msg))
	for i, v := range msg {
		str[i] = fmt.Sprintf("%02X", v)
	}
	return strings.Join(str, sep)
}
func byteToString(msg []byte, sep string) string {
	str := make([]string, len(msg))
	for i, v := range msg {
		str[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(str, sep)
}

//parsing chinese
func toUtf8(b []byte) string {
	b = GetValidByte(b)
	s, err := DecodeBig5(b)
	if err != nil {
		return string(b)
	}
	return string(s)
}

func DecodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
func GetValidByte(src []byte) []byte {
	var str_buf []byte
	for _, v := range src {
		if v != 0 {
			str_buf = append(str_buf, v)
		} else {

			break
		}
	}
	return str_buf
}

const ipmax = 250
const maxnumber = ipmax * ipmax //250*250

//Return TestParam
//
//Exampe:Hostname=device,number=1
//
//MACAddress: "00-60-E9-2C-00-"+nmuber=00-60-E9-2C-00-01
//
//IPAddress:"192.168.1." + number=192.168.1.1
//
//Netmask: "255.255.255.0"
//Gateway: "192.168.1.254"
//
//Hostname:Hostname + number =device1
func GetTestParam(Hostname string, number uint) (gwd.ModelInfo, error) {
	defaultv := 1
	if number > maxnumber {
		return gwd.ModelInfo{}, fmt.Errorf("number than %v", maxnumber)
	}
	if number == 0 {
		return gwd.ModelInfo{}, fmt.Errorf("number can't < %v", 0)
	}

	locfour := int(number) % ipmax
	var locthird int
	if locfour == 0 {
		locthird = int(number)/ipmax - 1
		locfour = ipmax

	} else {
		locthird = int(number) / ipmax
	}
	mac := "00-60-E9-2C-" + fmt.Sprintf("%02X", locthird) + "-" + fmt.Sprintf("%02X", locfour)
	ip := "192.168." + strconv.Itoa(locthird+defaultv) + "." + strconv.Itoa(locfour)
	gate := "192.168." + strconv.Itoa(locthird+defaultv) + "." + "254"
	n := strconv.Itoa(int(number))
	return gwd.ModelInfo{Model: "testdevice" + n,
		MACAddress: mac,
		IPAddress:  ip,
		Netmask:    "255.255.255.0",
		Gateway:    gate, Hostname: Hostname + n}, nil
}

//Return HostName
//
//Exampe:Hostname=device,number=1
//
//return device1
func GetTestHostName(Hostname string, number int) string {

	return Hostname + strconv.Itoa(number)
}

//Return IP
//
//Exampe:number=1
//
//return "192.168.30.1"
func GetTestNewIP(number int) string {

	return "192.168.30." + strconv.Itoa(number)
}

//Return IP
//
//Exampe:number=1
//
//return "255.255.255.1"
func GetTestNetMask(number int) string {

	return "255.255.255." + strconv.Itoa(number)
}

//Return IP
//
//Exampe:number=1
//
//return "192.168.255.1"
func GetGateway(number int) string {

	return "192.168.255." + strconv.Itoa(number)
}
