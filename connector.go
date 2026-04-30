package main

import (
	"crusher/crusher"
	"crusher/encryptor"
	"crusher/randomizer"
	"crusher/wrapper"
	"fmt"
	"io"
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
		f, err := os.Open("test/data.txt")
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
		defer wg.Done()
		crumbs := <-chanCrush
		for i, crumb := range crumbs {
			time.Sleep(time.Millisecond * 15)
			fmt.Println(i, " crumb --->", string(wrapper.Wrap(crumb)[36:36+len(crumb.Payload)]))
		}
		fmt.Println("Wrap - OK")
	}()
	wg.Wait()
}
