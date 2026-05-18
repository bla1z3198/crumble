package main

import (
	"crusher/encryptor"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

var (
	crmb      chan *Crumb
	done      chan []byte
	collected chan []byte
	flows     map[uint16]chan LCrumb
	client    *net.UDPAddr
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

type Cfg struct {
	Mtu  uint16 `json:"mtu"`
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
	// Init config
	cfg := Init()
	// Resolve and listen udp...
	addr, _ := net.ResolveUDPAddr("udp", cfg.IP+cfg.Port)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("server can't listen", addr.String(), err)
		return
	}
	// Close on interrupt
	defer conn.Close()
	// Success start messages
	fmt.Println("crumble server on", cfg.IP+cfg.Port, "...")
	fmt.Println("mtu =", cfg.Mtu)
	// Start multiplexer
	go Mux()
	// Define pack []byte array
	pack := make([]byte, cfg.Mtu)
	// Read from udp, decode, unwrap
	for {
		volume, user_addr, _ := conn.ReadFromUDP(pack)
		client = user_addr
		decode := encryptor.Encrypt(pack[:volume])
		Unwrap(decode[25:])
	}
}

func Mux() {
	for {
		// 1 case: crmb chan have data
		select {
		// Read from crmb
		case crumb := <-crmb:
			// Check records in flows map
			_, exist := flows[crumb.FlowID]
			// Create new record in flows map, if ID not exist
			if !exist {
				flows[crumb.FlowID] = make(chan LCrumb)
				go Flow(crumb.FlowID, flows[crumb.FlowID])
			}
			// Write crumb info to LCrumb chan
			flows[crumb.FlowID] <- LCrumb{
				crumb.Seq,
				crumb.Lost,
				crumb.Payload,
			}
		// 2 case: done chan have data
		case ready := <-done:
			// Unwrap id and payload
			id := binary.BigEndian.Uint16(ready[0:2])
			payload := ready[3:]
			go Duplex(client.String(), payload)
			// Close and delete if chan is available
			if _, ok := flows[id]; ok {
				close(flows[id])
				delete(flows, id)
			}
		}
	}
}

func Duplex(addr string, payload []byte) {
	// Dial to a client
	conn, err := net.Dial("udp", "192.168.1.153:9443")
	if err != nil {
		panic(err)
	}
	// Send payload to client
	conn.Write(payload)
}

func Flow(ID uint16, ch chan LCrumb) {
	// Init
	buf := make([]byte, 2)
	ready := make(map[uint16][]byte, 0)
	result := make([]byte, 0)
	for {
		// Read packets from ch
		pack, ok := <-ch
		if !ok {
			break
		}
		// Duplicate protection
		if _, exist := ready[pack.Seq]; exist {
			continue
			// Ordering packets in a correct sequence
		} else {
			ready[pack.Seq] = pack.Payload
		}
		// Condition for end of flow
		if len(ready) == int(pack.Lost) {
			for i := 0; i < len(ready); i++ {
				result = append(result, ready[uint16(i)]...)
			}
			fmt.Println("ID", ID, "LEN ->", len(result))
			// ID and payload in []byte array
			binary.BigEndian.PutUint16(buf[0:2], ID)
			buf = append(buf, result...)
			done <- buf
			break
		}
	}
}

func Unwrap(crumb []byte) {
	// Unwrap []byte into crumbs
	unpacked = &Crumb{
		FlowID:  uint16(crumb[0])<<8 | uint16(crumb[1]),
		Seq:     uint16(crumb[2])<<8 | uint16(crumb[3]),
		Flags:   string(crumb[4:8]),
		Lost:    uint16(crumb[9])<<8 | uint16(crumb[10]),
		Payload: crumb[11:],
		Padding: nil,
	}
	crmb <- unpacked
}

func Init() Cfg {
	// Read file and init config
	data, err := os.ReadFile("collector.json")
	if err != nil {
		panic(err)
	}
	cfg := Cfg{}
	json.Unmarshal(data, &cfg)

	collected = make(chan []byte, cfg.Mtu)
	crmb = make(chan *Crumb, cfg.Buf)
	flows = make(map[uint16]chan LCrumb, cfg.Buf)
	done = make(chan []byte, cfg.Buf)

	return cfg
}
