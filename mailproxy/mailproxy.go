package mailproxy

import (
	"bytes"

	"github.com/jkak/mail"
	rotlog "github.com/jkak/rotlogs/daterot"
)

type formMapping map[string]string
type fileMapping map[string]*bytes.Buffer

type chanMapping map[string]chan *mail.Message
type senderMapping map[string][]Sender

var (
	cfg       = new(Config)
	mailCh    = make(chanMapping)
	senderMap = make(senderMapping)

	Logger = rotlog.LoggerPtr.Logger
)

func Init(cfgFile string) {
	// read and decode config
	decodeConfigAndInit(cfgFile)
}
