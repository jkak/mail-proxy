package mailproxy

import (
	"bytes"
	"net/http"

	rotlog "github.com/jkak/rotlogs/daterot"
)

type formMapping map[string]string
type fileMapping map[string]*bytes.Buffer

type chanMapping map[string]chan *sendMsg
type senderMapping map[string][]Sender

var (
	cfg       = new(Config)
	mailCh    = make(chanMapping)
	senderMap = make(senderMapping)
	// Logger for log
	Logger = rotlog.LoggerPtr.Logger
)

// Init _
func Init(cfgFile string) {
	// read and decode config
	decodeConfigAndInit(cfgFile)
}

// Run _
func Run() {
	// run sender
	mailSenders()

	// bind handler
	http.HandleFunc(cfg.URI, handleMail)
	Logger.Fatal(http.ListenAndServe(cfg.ProxyPort, nil))
}
