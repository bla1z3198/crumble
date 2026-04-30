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
	for {
		f, err := os.Open("data.txt")
		if err != nil {
			fmt.Println("can't open data")
		}
		data, _ := io.ReadAll(f)
		f.Close()

		wg.Add(1)
		go func(data *[]byte) {
			parts, one, _ := randomizer.Random(len(*data))
			encrypted := encryptor.Encryption(*data)

			info := crusher.Service{
				Encrypted: encrypted,
				ID:        uint16(950),
				Flg:       "DATA",
				Parts:     uint16(parts),
				One:       uint16(one),
			}

			crumbs := crusher.Crush(&info)
			for i := range crumbs {
				wrapped := wrapper.Wrap(crumbs[i])
				time.Sleep(time.Millisecond * 10)
				fmt.Println(string(wrapped[36:]))
			}
			wg.Done()
		}(&data)
		wg.Wait()
	}
}

func Push(ready []byte, len int) {
}
