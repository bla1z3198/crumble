package main

import (
	"crusher/encryptor"
	"fmt"
	"net"
)

var (
	mtu       uint16
	ip        string
	port      string
	addr      *net.UDPAddr
	pipe      chan *Crumb
	done      chan uint16
	completed chan []byte
	flows     map[uint16]chan LCrumb
	pack      []byte
	unpacked  *Crumb
)

type Crumb struct {
	FlowID  uint16
	Seq     uint16
	Flags   string
	Lost    uint16
	Payload []byte
	Padding []byte
}

type Config struct {
	MTU  uint16 `json:"mtu"`
	IP   string `json:"ip"`
	Port string `json:"port"`
	Buf  uint16 `json:"buf"`
}

type LCrumb struct {
	Seq     uint16
	Lost    uint16
	Payload []byte
}

func main() {
	ip, port, mtu, addr = Init()
	fmt.Println("mtu =", mtu)

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("server can't listen", ip+port, err)
		return
	}
	defer conn.Close()

	fmt.Println("crumble server on", ip+port, "...")

	go Mux()
	go Tun()

	for {
		volume, user_addr, _ := conn.ReadFromUDP(pack)
		fmt.Println("got pack from:", user_addr.String())
		decode := encryptor.Encrypt(pack[:volume])
		Unwrap(decode[25:])
	}
}

func Mux() {
	for {
		select {
		case crumb := <-pipe:
			_, exist := flows[crumb.FlowID]

			if !exist {
				flows[crumb.FlowID] = make(chan LCrumb)
				go Flow(crumb.FlowID, flows[crumb.FlowID])
			}

			flows[crumb.FlowID] <- LCrumb{
				crumb.Seq,
				crumb.Lost,
				crumb.Payload,
			}

		case id := <-done:
			_, ok := flows[id]
			if ok {
				close(flows[id])
				delete(flows, id)
			}
		}
	}
}

func Tun() {
	conn, _ := net.Dial("udp", "192.168.1.153:888")
	fmt.Println("connected to client....")
	for {
		duplex := <-completed
		answer := []byte("hey connector, i am shaking your hand!")
		answer = append(answer, duplex...)
		conn.Write(answer)
	}
}

func Flow(ID uint16, ch chan LCrumb) {
	ready := make(map[uint16][]byte, 0)
	result := make([]byte, 0)
	for {
		pack, ok := <-ch
		if !ok {
			break
		}
		if _, exist := ready[pack.Seq]; exist {
			continue
		} else {
			ready[pack.Seq] = pack.Payload
		}
		if len(ready) == int(pack.Lost) {
			for i := 0; i < len(ready); i++ {
				result = append(result, ready[uint16(i)]...)
			}
			fmt.Println("ID", ID, "LEN ->", len(result))
			completed <- result
			done <- ID
			break
		}
	}
}

func Unwrap(crumb []byte) {
	unpacked = &Crumb{
		FlowID:  uint16(crumb[0])<<8 | uint16(crumb[1]),
		Seq:     uint16(crumb[2])<<8 | uint16(crumb[3]),
		Flags:   string(crumb[4:8]),
		Lost:    uint16(crumb[9])<<8 | uint16(crumb[10]),
		Payload: crumb[11:],
		Padding: nil,
	}
	pipe <- unpacked
}

func Init() (string, string, uint16, *net.UDPAddr) {
	cfg := Config{
		1280,
		"192.168.1.11",
		":999",
		1024,
	}

	addr, _ := net.ResolveUDPAddr("udp", cfg.IP+cfg.Port)
	pack = make([]byte, cfg.MTU)
	completed = make(chan []byte, cfg.MTU)
	pipe = make(chan *Crumb, cfg.Buf)
	flows = make(map[uint16]chan LCrumb, cfg.Buf)
	done = make(chan uint16, cfg.Buf)

	return cfg.IP, cfg.Port, cfg.MTU, addr
}
