package main

import (
	"flag"

	mp "github.com/jkak/mail-proxy/mailproxy"
)

var cfgFile string

func main() {
	flag.StringVar(&cfgFile, "c", "config.toml", "mail config file")
	flag.Parse()

	mp.Init(cfgFile)
	logger := mp.Logger
	logger.Infof("main() starting")

	mp.Run()
}
