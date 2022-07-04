package FirmWare

func firmWarePacket() []byte {
	packet := make([]byte, 40)
	def := "name1234passwd12modelname 123456" //not important, just input anychar
	for i, v := range def {
		packet[i] = byte(v)
	}
	packet[36] = 0x72
	return packet
}
