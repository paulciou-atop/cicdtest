package gwd

func invitePacket() []byte {

	packet := make([]byte, 300)
	packet[0] = 2
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet

}

func configPacket() []byte {

	packet := make([]byte, 300)
	packet[0] = 0
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet

}
func rebootPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 5
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet
}

func beepPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 7
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet
}

func reSetDefaultPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 5
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet
}
