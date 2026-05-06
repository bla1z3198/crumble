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
	chanCrush := make(chan []crusher.Crumb)
	wg.Add(3)
	go func() {
		defer wg.Done()
		f, err := os.Open("tests/plain_text.txt")
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
		data := <-chanRaw
		parts, one := randomizer.Random(len(data))
		info := crusher.Service{
			Encrypted: data,
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
			jitter := time.Millisecond * 4
			time.Sleep(jitter)
			ready := encryptor.Encrypt(wrapper.Wrap(crumb))
			conn.Write(ready)
		}
		fmt.Println("Wrapped and send!")
	}()
	wg.Wait()
}
