package main

import (
	"fmt"
	"net"
)

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

		ans := string(buf[36:n])
		fmt.Println("Payload:", ans)
	}
}
