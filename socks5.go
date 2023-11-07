package main

import (
	"fmt"

	"github.com/armon/go-socks5"
)

type socks5_credStruct struct {
	username string
	password string
}

var socks5_creds socks5_credStruct

func socks5_init() {
	defer func() {
		recover()
	}()

	socks5_creds.username = shortID(8)
	socks5_creds.password = shortID(8)

	fmt.Printf("Socks initialized: %s:%s\n", socks5_creds.username, socks5_creds.password)

	// Create a SOCKS5 server
	conf := &socks5.Config{Credentials: socks5.StaticCredentials{socks5_creds.username: socks5_creds.password}}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8997
	if err := server.ListenAndServe("tcp", "0.0.0.0:8997"); err != nil {
		fmt.Println("Unable to open socks server")
		return
	}
}

func init_socks5() {
	go socks5_init()
}
