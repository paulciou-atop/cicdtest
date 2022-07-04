package gwd

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

//paring ip  Split "."
func parsingIp(ip string) ([]byte, error) {
	v := strings.Split(ip, ".")
	if len(v) != 4 {
		return nil, errors.New("IPAddress format error")
	}
	var ips [4]byte
	for i := 0; i < 4; i++ {
		r, err := strconv.Atoi(v[i])
		if err != nil {
			return nil, fmt.Errorf("IPAddress format error ,reason:%s", err.Error())
		}
		ips[i] = byte(r)
	}

	return ips[:], nil
}

//paring MacAddress  Split ":"
func parsingMacAddress(mac string) ([]byte, error) {
	m := strings.ReplaceAll(mac, "-", "")
	v, err := hex.DecodeString(m)
	if err != nil {
		return nil, fmt.Errorf("MacAddress format error ,reason:%s", err.Error())
	}
	return v, nil
}

func EncodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
