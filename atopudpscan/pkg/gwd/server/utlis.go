package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
