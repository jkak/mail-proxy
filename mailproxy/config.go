package mailproxy

import (
	"log"

	"github.com/BurntSushi/toml"
	rotlog "github.com/jkak/rotlogs/daterot"
)

// Config for read from conf file
type Config struct {
	ServerHost string // mail server host
	ServerPort int    // mail server port

	ProxyPort      string // proxy port for local
	ProxyKeepalive int    // proxy server keepalive a connection within seconds
	RetryInterval  int    // proxy server sleep second before retry
	RetryTimes     int    // proxy server retry times when send error
	Timeout        int    // timeout in seconds for internal process

	StatusOK string // status ok for response
	StatusEr string // status error for response

	LogFile string   // log file name
	URI     string   // uri for client
	Debug   bool     // log debug switch
	ContLen int      // mail content length
	Senders []Sender // sender info for send mail
}

// Sender struct for sender info
type Sender struct {
	SendName string // sender name
	MailUser string // who send the email
	Password string // password of whom
	Length   int    // sender channel length
}

func decodeConfigAndInit(f string) {
	_, err := toml.DecodeFile(f, cfg)
	if err != nil {
		log.Fatalf("decode cfg file err:%s\n", err)
	}

	// setting Logger
	rotlog.BaseFileName = cfg.LogFile
	rotlog.BaseLinkName = cfg.LogFile
	if cfg.Debug {
		rotlog.LogLevel = rotlog.DebugLevel
	} else {
		rotlog.LogLevel = rotlog.InfoLevel
	}
	Logger, err = rotlog.Rotate()
	if err != nil {
		log.Fatalf("create rotate log file err:%s", err)
	}

	// init senderMap
	for _, send := range cfg.Senders {
		senderMap[send.SendName] = append(senderMap[send.SendName], send)
		Logger.Printf("sender: %s, user: %s", send.SendName, send.MailUser)
	}

	// init mail channel
	for _, sender := range cfg.Senders {
		mailCh[sender.MailUser] = make(chan *sendMsg, sender.Length)
	}

	for k, v := range senderMap {
		Logger.Printf("sender: %s, Len:%d", k, len(v))
	}
}
