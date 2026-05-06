package main

import (
	"crusher/encryptor"
	"fmt"
	"net"
)

var (
	unpacked *Crumb
)

type Crumb struct {
	FlowID  uint16
	Seq     uint16
	Flags   string
	Lost    uint16
	Payload string
	Padding []byte
}

func main() {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:5252")

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Err!", err)
		return
	}
	defer conn.Close()

	fmt.Println("crumble on 127.0.0.1:5252...")

	buf := make([]byte, 1280)
	for {
		n, _, _ := conn.ReadFromUDP(buf)
		u := unpacked.Unwrap(encryptor.Encrypt(buf[:n])[25:])

		fmt.Printf("FlowID: %d\nSeq: %d\nFlags: %s\nLost: %d\nPayload: %s\n",
			u.FlowID,
			u.Seq,
			u.Flags,
			u.Lost,
			u.Payload)
		fmt.Println("<---------------------------------------->")
	}
}

func (unwrapped *Crumb) Unwrap(crumb []byte) *Crumb {
	unwrapped = &Crumb{
		FlowID:  0,
		Seq:     0,
		Flags:   "",
		Lost:    0,
		Payload: "",
		Padding: nil,
	}
	unwrapped.FlowID = uint16(crumb[0])<<8 | uint16(crumb[1])
	unwrapped.Seq = uint16(crumb[2])<<8 | uint16(crumb[3])
	unwrapped.Flags = string(crumb[4:8])
	unwrapped.Lost = uint16(crumb[9])<<8 | uint16(crumb[10])
	unwrapped.Payload = string(crumb[11:])

	return unwrapped
}
