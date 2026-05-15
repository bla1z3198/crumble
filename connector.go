package main

import (
	"crusher/crusher"
	"crusher/encryptor"
	"crusher/randomizer"
	"crusher/wrapper"
	"fmt"
	"net"
	"sync"

	"github.com/songgao/water"
)

var (
	wg    sync.WaitGroup
	mu    sync.Mutex
	tun   chan []byte
	conn2 *net.UDPConn
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", "192.168.1.153:888")
	conn2, _ = net.ListenUDP("udp", addr)

	tun = make(chan []byte, 128)
	conn, _ := net.Dial("udp", "192.168.1.11:999")
	wg.Add(100)

	go Tun()
	go Listen()

	for i := 0; i < 100; i++ {
		go PipeLine(i, conn)
	}
	wg.Wait()
}

func PipeLine(i int, conn net.Conn) {
	defer wg.Done()
	data := <-tun

	mu.Lock()
	parts, one := randomizer.Random(len(data))
	mu.Unlock()

	mu.Lock()
	info := crusher.Service{
		Encrypted: data,
		ID:        uint16(i),
		Flg:       "DATA",
		Parts:     uint16(parts),
		One:       uint16(one),
	}
	mu.Unlock()

	mu.Lock()
	crumbs := crusher.Crush(&info)
	mu.Unlock()

	for _, crumb := range crumbs {
		//jitter := time.Millisecond * 4
		//time.Sleep(jitter)
		mu.Lock()
		wrapped := wrapper.Wrap(crumb)
		ready := encryptor.Encrypt(wrapped)
		conn.Write(ready)
		mu.Unlock()
	}
}

func Tun() {
	config := water.Config{
		DeviceType: water.TUN,
	}

	ifce, err := water.New(config)
	if err != nil {
		panic(err)
	}

	fmt.Println("TUN:", ifce.Name())

	buf := make([]byte, 1280)

	for {
		n, _ := ifce.Read(buf)
		pkt := buf[:n]
		tun <- pkt
	}
}

func Listen() {
	buf := make([]byte, 1280)
	for {
		volume, _, _ := conn2.ReadFromUDP(buf)
		fmt.Println(string(buf[:volume]))
	}
}
