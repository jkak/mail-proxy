package main

import (
	"flag"
	"log"
	"net/http"

	rotlog "github.com/jkak/rotlogs/daterot"
)

var (
	cfgf   string
	err    error
	logger = rotlog.LoggerPtr.Logger
)

func init() {
	flag.StringVar(&cfgf, "c", "test.toml", "mail config file")
	flag.Parse()

	// read and decode config
	decodeConfig()

	// setting logger
	rotlog.BaseFileName = cfg.LogFile
	if cfg.Debug {
		rotlog.LogLevel = rotlog.DebugLevel
	} else {
		rotlog.LogLevel = rotlog.InfoLevel
	}
	logger, err = rotlog.Rotate()
	if err != nil {
		log.Fatalf("create rotate log file err:%s\n", err)
	}
}

func main() {
	logger.Infof("main sender num:%+v", len(senderMap))

	http.HandleFunc("/mail", handleMail)
	logger.Fatal(http.ListenAndServe(cfg.ProxyPort, nil))
}
