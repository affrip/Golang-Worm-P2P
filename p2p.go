package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var peers []string

func init_peers() {
	peers = append(peers, GetOutboundIP().String())
	pf, err := os.Open(".book.dat")

	if err != nil {
		panic("Cannot open book.dat")
	}

	scanner := bufio.NewScanner(pf)

	for scanner.Scan() {
		peers = append(peers, strings.TrimSpace(scanner.Text()))
	}

	pf.Close()
}

func peer_exists(peer string) bool {
	for i := 0; i < len(peers); i++ {
		if peer == peers[i] {
			return true
		}
	}
	return false
}

func test_peer(peer string) bool {
	fmt.Println("Trying to dial peer: ", peer)

	conn, err := net.Dial("udp", peer+":33069")
	if err != nil {
		fmt.Println("Error dialing peer")
		return false
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(3 * time.Second))

	_, err = fmt.Fprintf(conn, "ping\x00")
	if err != nil {
		fmt.Println("Error writing to peer: ", peer)
		return false
	}

	connbuf := bufio.NewReader(conn)

	data, err := connbuf.ReadString('\x00')
	if err != nil {
		return false
	}

	if data == "" {
		fmt.Println("Didn't get response, ", peer, " offline.")
		return false
	}

	if strings.Split(data, "\x00")[0] == "pong" {
		fmt.Println("Peer speaks soy ", peer)
		return true
	}

	return false
}

func save_peer(peer string) {
	if test_peer(peer) == false {
		return
	}

	peers = append(peers, peer)
	pf, err := os.OpenFile(".book.dat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error opening peer db")
		return
	}

	pf.Write([]byte(peer))
	pf.Write([]byte("\n"))

	pf.Close()
}

func crawlpeer(peer string) {
	fmt.Println("Trying to dial peer: ", peer)

	conn, err := net.Dial("udp", peer+":33069")
	if err != nil {
		fmt.Println("Error dialing peer")
		return
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(3 * time.Second))

	_, err = fmt.Fprintf(conn, "getpeers\x00")
	if err != nil {
		fmt.Println("Error writing to peer: ", peer)
		return
	}

	connbuf := bufio.NewReader(conn)
	data, err := connbuf.ReadString('\x00')
	if err != nil {
		return
	}

	if data == "" {
		fmt.Println("Didn't get response, retrying peer ", peer, " later.")
		return
	}

	data = strings.TrimSpace(strings.Split(data, "\x00")[0])
	fmt.Println("Got peers from: ", peer)

	peer_arr := strings.Split(data, "\n")

	for i := 0; i < len(peer_arr); i++ {
		if peer_exists(peer_arr[i]) == false {
			save_peer(peer_arr[i])
			fmt.Println("Saved peer: ", peer_arr[i])
		} else {
			fmt.Println("Peer already exists: ", peer_arr[i])
		}
	}
}

func push_config(peer string) {
	if encrypted_config == "" {
		fmt.Println("Skipping push config not received valid config yet")
		return
	}

	fmt.Println("Trying to dial peer: ", peer)

	conn, err := net.Dial("udp", peer+":33069")
	if err != nil {
		fmt.Println("Error dialing peer")
		return
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(3 * time.Second))

	_, err = fmt.Fprintf(conn, "pushconfig\x00%s\x00", encrypted_config)
	if err != nil {
		fmt.Println("Error writing to peer: ", peer)
		return
	}
}

func pull_config(peer string) {
	fmt.Println("Trying to dial peer: ", peer)

	conn, err := net.Dial("udp", peer+":33069")
	if err != nil {
		fmt.Println("Error dialing peer")
		return
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(3 * time.Second))

	_, err = fmt.Fprintf(conn, "pullconfig\x00")
	if err != nil {
		fmt.Println("Error writing to peer: ", peer)
		return
	}

	connbuf := bufio.NewReader(conn)
	data, err := connbuf.ReadString('\x00')
	if err != nil {
		return
	}

	if data == "" {
		fmt.Println("Didn't get response, retrying peer ", peer, " later.")
		return
	}

	data = strings.TrimSpace(strings.Split(data, "\x00")[0])
	fmt.Println("Got config from: ", peer, "\n", data)
}

func p2p_client() {
	for {
		shufarray(peers)

		for i := 0; i < len(peers); i++ {
			crawlpeer(peers[i])
			time.Sleep(1 * time.Second)
		}

		if developer_mode == false {
			time.Sleep(60 * time.Second)
		}

		for i := 0; i < len(peers); i++ {
			push_config(peers[i])
			time.Sleep(1 * time.Second)
		}

		if developer_mode == false {
			time.Sleep(60 * time.Second)
		}

		for i := 0; i < len(peers); i++ {
			pull_config(peers[i])
			time.Sleep(1 * time.Second)
		}

		if developer_mode == false {
			time.Sleep(60 * time.Second)
		}
	}
}

func get_peerlist() (peerlist string) {
	return strings.Join(peers, "\n")
}

func get_pushed_config(commandbody string) {
	if verify_message_json([]byte(commandbody)) == nil {
		var commandjson messageStruct

		err := json.Unmarshal([]byte(commandbody), &commandjson)
		if err != nil {
			fmt.Println("Error unmarshaling message struct")
		}

		fmt.Println("Valid message received...")
		check_new_config([]byte(commandjson.Msg), commandbody)
	}
}

func sendResponse(conn *net.UDPConn, addr net.Addr, command string) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	commandbody := strings.Split(command, "\x00")[1]
	command = strings.Split(command, "\x00")[0]
	fmt.Println(command)

	switch command {
	case "getpeers":
		fmt.Println("Answering getpeers to: ", addr.String())
		conn.WriteTo([]byte(get_peerlist()+"\x00"), addr)
		break

	case "pullconfig":
		fmt.Println("Sending encrypted config to: ", addr.String())
		conn.WriteTo([]byte(encrypted_config+"\x00"), addr)
		break

	case "pushconfig":
		fmt.Println("Getting pushed config from: ", addr.String())
		get_pushed_config(commandbody)
		break

	default:
		conn.WriteTo([]byte("pong"+"\x00"), addr)
		break
	}
}

func p2p_server() {
	p := make([]byte, 4096)

	addr := net.UDPAddr{
		Port: 33069,
		IP:   net.ParseIP("0.0.0.0"),
	}

	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	defer ser.Close()

	for {
		ser.SetReadDeadline(time.Now().Add(5 * time.Second))
		fmt.Println("Reading from udp server...")
		_, remoteaddr, err := ser.ReadFrom(p)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("timeout")
				continue
			}
			fmt.Println("Error reading from udp ", err.Error())
			continue
		}

		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)

		remoteip := remoteaddr.String()

		if peer_exists(remoteip) == false {
			save_peer(remoteip)
			fmt.Println("Saved remote peer: ", remoteip)
		} else {
			fmt.Println("Remote peer already exists: ", remoteip)
		}

		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

		go sendResponse(ser, remoteaddr, string(p))
	}
}

func init_p2p() {
	init_peers()
	go p2p_server()

	time.Sleep(1 * time.Second)
	go p2p_client()
}
