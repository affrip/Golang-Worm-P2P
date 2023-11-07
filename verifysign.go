package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
)

type messageStruct struct {
	Msg        []byte `json:"msg"`
	MsgHashSum []byte `json:"msghashnum"`
	Signature  []byte `json:"signature"`
}

var public_key = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp4JiM6Bq3Yh24+JGC8MdGk/QoQhxWigYJJrsCstmiJNJCWpSuxk+xbdQIgv1E29sJ7GQ4skOtUzoj7dS5QtmgNcPrlOt7xzS4qEfiuGnInraXP4XPaXfOsU+Br9Oy2VsVirkH+hcZ11bWvWblavSM7iqMukF1KGySrOlLRTtUyRhQEUJWRBsNhnZTa0SlHKfXNGjWERMTaGQQ60eu/1NKEE2iVkoifuuQvuMnOcOt96m1EjsgVdG05/4JHiwi+ojBOpkvfgAgb9KYaTNYpBYBdtmGFNpdERt5Z6ezjx7YFoXctlVoZ/7SzMmadOiEJn8BmYzXs5Co/l7hZ7fHhUXUwIDAQAB\n-----END PUBLIC KEY-----\n"

func genkeypair() (rsa.PrivateKey, rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return rsa.PrivateKey{}, rsa.PublicKey{}, errors.New("Error generating key pair")
	}
	// The public key is a part of the *rsa.PrivateKey struct
	publicKey := privateKey.PublicKey
	return *privateKey, publicKey, nil
}

func verify_message(message []byte, msgHashSum []byte, signature []byte) error {
	pkey, err := ParseRsaPublicKeyFromPemStr(public_key)
	if err != nil {
		return errors.New("Error loading public key")
	}

	err = rsa.VerifyPSS(pkey, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return errors.New("Could not verify signature")
	}

	fmt.Println(msgHashSum)
	fmt.Println("Message verified: " + string(message))

	return nil
}

func verify_message_json(jsdata []byte) error {
	var messageObj messageStruct

	err := json.Unmarshal(jsdata, &messageObj)

	if err != nil {
		fmt.Println("Error unmarshalling message", string(jsdata))
		return err
	}

	err = verify_message(messageObj.Msg, messageObj.MsgHashSum, messageObj.Signature)
	return err
}
