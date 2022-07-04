package FirmWare

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const fwport = 55950
const size = 512
const upgradetimeout = 180 * time.Second
const resposetimeout = 5 * time.Second
const conntimeout = 10 * time.Second

type fwStatus = string
type ProcessStatus = string

const (
	None      ProcessStatus = "none"
	Upgrading ProcessStatus = "Upgrading"
	Complete  ProcessStatus = "complete"
	Error     ProcessStatus = "Error"
)

const (
	ready  fwStatus = "E001"
	erased fwStatus = "S001"
	finish fwStatus = "S002"
	going  fwStatus = "a"
)

func NewFirmWare(ip string) *FirmWare {
	return &FirmWare{ip: ip, m: new(sync.Mutex), process: None, c: make(chan bool, 1)}
}

type FirmWare struct {
	ip           string
	filesize     int64
	process      string
	m            *sync.Mutex
	c            chan bool
	errorMessage string
}

func (f *FirmWare) GetProcessStatus() (string, error) {
	f.m.Lock()
	s := f.process
	m := f.errorMessage
	f.m.Unlock()
	if m == "" {
		return s, nil
	}
	return s, errors.New(m)
}
func (f *FirmWare) Init() {
	f.settingProcess(None)
	f.writeLog("")
	f.c = make(chan bool, 1)

}
func (f *FirmWare) Upgrading(file *os.File) error {
	f.Init()
	packetCount := 0
	fwHeader := false
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	f.filesize = fi.Size()
	address := strings.Join([]string{f.ip, strconv.Itoa(fwport)}, ":")
	/*=== check device fw port ====*/
	ckcon, err := net.DialTimeout("tcp", address, conntimeout)
	if err != nil {
		f.settingProcess(Error)
		return fmt.Errorf("ip:%v can't connect", f.ip)
	}

	ckcon.Close()
	/*=== check device fw port done ====*/

	/*=== update fw ====*/
	go func() {
		f.c <- true
		upconn, err := net.DialTimeout("tcp", address, conntimeout)
		if err != nil {
			f.writeLog(err.Error())
			f.settingProcess(Error)
			return
		}

		defer upconn.Close()

		buf := make([]byte, 0, size)
		r := bufio.NewReader(file)
		for {
			//send fw header packet
			if !fwHeader {
				_, err := upconn.Write(downloadRequest(f.filesize))
				if err != nil {
					f.settingProcess(Error)
					f.writeLog(err.Error())
					return
				}

				err = f.waitResponse(upconn, going)
				if err != nil {
					f.settingProcess(Error)
					f.writeLog(err.Error())
					return
				}

				fwHeader = true
				f.calculateProcess(packetCount)

			}

			n, readerr := io.ReadFull(r, buf[:cap(buf)])
			packetCount++
			f.calculateProcess(packetCount)

			buf = buf[:n]
			_, err := upconn.Write(buf)
			if err != nil {
				f.settingProcess(Error)
				f.writeLog(err.Error())
				return
			}

			err = f.waitResponse(upconn, going)
			if err != nil {
				f.settingProcess(Error)
				f.writeLog(err.Error())
				return
			}
			/*=== wait updraging fw ====*/
			if readerr != nil {
				if readerr == io.EOF || readerr == io.ErrUnexpectedEOF {

					err := f.waitResponse(upconn, erased)
					if err != nil {
						f.settingProcess(Error)
						f.writeLog(err.Error())
						return
					}
					break
				}
			}
			/*=== wait updraging fw done====*/
		}
		/*=== update fw done ====*/

		/*=== wait finish ====*/
		err = f.waitResponse(upconn, finish)
		if err != nil {
			f.settingProcess(Error)
			f.writeLog(err.Error())
			return
		}
	}()
	<-f.c
	/*=== wait finish done ====*/
	return nil
}

func (f *FirmWare) waitResponse(con net.Conn, w fwStatus) error {

	if w == going {
		con.SetReadDeadline(time.Now().Add(resposetimeout))
	} else {
		con.SetReadDeadline(time.Now().Add(upgradetimeout))
	}

	dst := make([]byte, len(w))
	for {
		_, err := con.Read(dst)
		if err != nil {
			return err
		}
		r := strings.TrimSpace(string(dst))
		if r == w {
			if w == erased {
				f.settingProcess(Upgrading)
			}
			if w == finish {
				f.settingProcess(Complete)
			}
			return nil
		}

	}
}
func (f *FirmWare) calculateProcess(packetCount int) {
	proc := packetCount * 100 / int(math.Ceil(float64(f.filesize)/512))
	progressPercent := int(math.Floor(float64(proc)*100) / 100)
	f.settingProcess(strconv.Itoa(progressPercent))
}

func (f *FirmWare) settingProcess(process string) {
	f.m.Lock()
	f.process = process
	f.m.Unlock()
}

func (f *FirmWare) writeLog(m string) {
	f.m.Lock()
	f.errorMessage = m
	f.m.Unlock()

}

func downloadRequest(filesize int64) []byte {
	dl_request := firmWarePacket()
	//dl_request[32] ~ dl_request[35] :save file size
	for j := 3; j >= 0; j-- {
		dl_request[j+32] = (byte)(filesize / int64(math.Pow(256, float64(j))))
		filesize = filesize - int64(dl_request[j+32])*int64(math.Pow(256, float64(j)))
	}
	return dl_request
}
