package main

import (
	"fmt"
	"log"

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
	ProxySleep int    // proxy server sleep second before retry
	ProxyPort  string // proxy port for local

	LogFile string   // log file name
	Debug   bool     // log debug switch
	ContLen int      // mail content length
	Senders []Sender // sender info for send mail
}

// Sender struct for sender info
type Sender struct {
	Name     string // sender name
	Account  string // who send the email
	Password string // password of whom
}

func decodeConfig() {
	_, err := toml.DecodeFile(cfgf, &cfg)
	if err != nil {
		log.Fatalf("decode cfg err:%s\n", err)
	}

	senderMap = make(map[string][]Sender)
	for _, send := range cfg.Senders {
		senderMap[send.Name] = append(senderMap[send.Name], send)
	}
	for k, v := range senderMap {
		fmt.Printf("sender: %s\tLen:%d\n", k, len(v))
	}
}
