package main

import (
	"crusher/crusher"
	"crusher/encryptor"
	"crusher/randomizer"
	"crusher/wrapper"
	"fmt"
	"net"
	"time"

	"github.com/songgao/water"
)

var (
	connL *net.UDPConn
	connD net.Conn
	tun   chan []byte
)

type Cfg struct {
	ip       string
	port     string
	password string
	mtu      int
}

func main() {
	// Init
	cfg := Init()
	fmt.Println("destination ip:port", cfg.ip+cfg.port)
	// Connect to the server
	connD, _ = net.Dial("udp", cfg.ip+cfg.port)
	fmt.Println("success, connected to:", cfg.ip+cfg.port)
	// Start TUN interface
	go Tun()
	// Start duplex link
	go Duplex()
	// Init Flow ID
	var ID uint16
	// For loop, go Flows
	for {
		data := <-tun
		go Flow(ID, data)
		ID++
		// Reset ID to zero
		if ID == 65534 {
			ID = 0
		}
	}
}

func Flow(id uint16, data []byte) {
	// Get info about divide
	rand := make(chan []int, 1)
	go randomizer.Random(len(data), rand)
	r := <-rand
	// Create struct about crushing, crush
	info := crusher.Service{
		Encrypted: data,
		ID:        id,
		Flg:       "DATA",
		Parts:     uint16(r[0]),
		One:       uint16(r[1]),
	}
	crumbs := make(chan []crusher.Crumb)
	go crusher.Crush(info, crumbs)
	c := <-crumbs
	// Wrap, encrypt, send
	for _, crumb := range c {
		jitter := time.Millisecond * 4
		time.Sleep(jitter)

		wrapped := wrapper.Wrap(crumb)
		ready := encryptor.Encrypt(wrapped)
		connD.Write(ready)
		fmt.Println("N ->", id)
	}
}

func Tun() {
	// TUN cfg
	config := water.Config{
		DeviceType: water.TUN,
	}
	// New TUN
	ifce, err := water.New(config)
	if err != nil {
		panic(err)
	}
	// TUN name
	fmt.Println("TUN:", ifce.Name())
	// Buffer for packets
	buf := make([]byte, 1280)
	// For loop: get data from TUN
	for {
		n, _ := ifce.Read(buf)
		pkt := buf[:n]
		tun <- pkt
	}
}

func Duplex() {
	// Resolve and listen
	self_ip := "192.168.1.153:9443"
	addr, _ := net.ResolveUDPAddr("udp", self_ip)
	connL, _ = net.ListenUDP("udp", addr)
	// Buffer for packets
	buf := make([]byte, 1280)
	// For loop: get loopback from server
	for {
		volume, _, _ := connL.ReadFromUDP(buf)
		fmt.Println(string(buf[:volume]))
	}
}

func Init() *Cfg {
	tun = make(chan []byte, 256)

	cfg := &Cfg{
		"192.168.1.11:",
		"9443",
		"crumble",
		1280,
	}
	return cfg
}
