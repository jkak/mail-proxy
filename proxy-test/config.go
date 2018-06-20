package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	cfg       Config
	senderMap map[string][]Sender
)

// Config for read from conf file
type Config struct {
	ServerHost string // mail server host
	ServerPort int    // mail server port
	ProxyPort  string // proxy port for local

	LogFile string   // log file name
	Debug   bool     // log debug switch
	ContLen int      // mail content length
	Senders []string // sender info for send mail
}

// Sender struct for sender info
type Sender struct {
	Account  string // who send the email
	Password string // password of whom
}

func decodeConfig() {
	_, err := toml.DecodeFile(cfgf, &cfg)
	if err != nil {
		log.Fatalf("decode cfg err:%s\n", err)
	}
	// fmt.Printf("cfg:%#v\n", cfg)

	senderMap = make(map[string][]Sender)
	for _, send := range cfg.Senders {
		// fmt.Println(send)
		sInfo := strings.Split(send, ":")
		sID := sInfo[0]
		sender := Sender{Account: sInfo[1], Password: sInfo[2]}
		senderMap[sID] = append(senderMap[sID], sender)
	}
	for k, v := range senderMap {
		fmt.Printf("sender: %s\tLen:%d\n", k, len(v))
	}
}
