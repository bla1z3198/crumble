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
	mu sync.Mutex
)

func main() {
	start := time.Now()
	conn, _ := net.Dial("udp", "127.0.0.1:5252")
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go PipeLine(i, conn)
	}
	wg.Wait()
	a := time.Since(start)
	fmt.Println("elapsed time ->", a)
}

func PipeLine(i int, conn net.Conn) {
	defer wg.Done()

	mu.Lock()
	f, err := os.Open("tests/plain_text.txt")
	if err != nil {
		fmt.Println("can't open data")
	}
	data, _ := io.ReadAll(f)
	f.Close()
	mu.Unlock()

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

	for i, crumb := range crumbs {
		//jitter := time.Millisecond * 4
		//time.Sleep(jitter)
		mu.Lock()
		wrapped := wrapper.Wrap(crumb)
		ready := encryptor.Encrypt(wrapped)
		fmt.Println("SHIPPED", 950+i)
		fmt.Println("PAYLOAD", string(wrapped[36:]))
		fmt.Println("Wrapped and send!", i)
		conn.Write(ready)
		mu.Unlock()
	}
}
