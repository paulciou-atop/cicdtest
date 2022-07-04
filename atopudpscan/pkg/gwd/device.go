package gwd

import (
	FirmWare "nms/atopudpscan/pkg/AtopFirmWare"
)

func NewDevice(ip string) *Device {
	return &Device{ip: ip}
}

type Device struct {
	ip string
	fw *FirmWare.FirmWare
}

func (d *Device) Ip() string {
	return d.ip
}

func (d *Device) FirmWare() *FirmWare.FirmWare {
	if d.fw == nil {
		d.fw = FirmWare.NewFirmWare(d.ip)
	}
	return d.fw

}
