package main

import (
	"crusher/crusher"
	"crusher/encryptor"
	"crusher/randomizer"
	"crusher/wrapper"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func main() {
	chanRaw := make(chan []byte, 1280)
	chanEncrypted := make(chan []byte, 1280)
	chanCrush := make(chan []crusher.Crumb)
	wg.Add(4)
	go func() {
		defer wg.Done()
		f, err := os.Open("tests/special_symbols.txt")
		if err != nil {
			fmt.Println("can't open data")
		}
		data, _ := io.ReadAll(f)
		f.Close()

		chanRaw <- data
		fmt.Println("Raw - OK")
	}()

	go func() {
		defer wg.Done()
		encrypted := encryptor.Encryption(<-chanRaw)
		chanEncrypted <- encrypted
		fmt.Println("Encrypt - OK")
	}()

	go func() {
		defer wg.Done()
		encrypted := <-chanEncrypted
		parts, one := randomizer.Random(len(encrypted))
		info := crusher.Service{
			Encrypted: encrypted,
			ID:        uint16(950),
			Flg:       "DATA",
			Parts:     uint16(parts),
			One:       uint16(one),
		}
		chanCrush <- crusher.Crush(&info)
	}()

	go func() {
		conn, _ := net.Dial("udp", "127.0.0.1:5252")
		defer conn.Close()
		defer wg.Done()
		crumbs := <-chanCrush
		for _, crumb := range crumbs {
			jitter := time.Millisecond * 15
			time.Sleep(jitter)
			ready := wrapper.Wrap(crumb)
			conn.Write(ready)
		}
		fmt.Println("Wrapped and send!")
	}()
	wg.Wait()
}
