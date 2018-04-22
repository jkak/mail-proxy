// https://github.com/jordan-wright/email
// https://godoc.org/github.com/jordan-wright/email
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/jordan-wright/email"
)

var (
	cfgf string
	cfg  Config
	err  error
)

// Config for read from conf file
type Config struct {
	ServerHost string // mail server host
	ServerPort int    // mail server port

	FromUser string // who send the email
	FromNick string // nick name of who
	Password string // password of who
}

func init() {
	flag.StringVar(&cfgf, "c", "test.toml", "mail config file")
	flag.Parse()

	// read config
	_, err := toml.DecodeFile(cfgf, &cfg)
	if err != nil {
		log.Fatalf("decode cfg err:%s\n", err)
	}
	log.Printf("cfg:%+v\n", cfg)
}

func handleMail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	content := r.FormValue("content")
	subject := r.FormValue("subject")
	sender := r.FormValue("sender")
	tos := r.FormValue("tos")

	if len(content) == 0 || len(subject) == 0 ||
		len(tos) == 0 || len(sender) == 0 {
		w.Write([]byte(`err: some paramter is blank`))
		return
	}
	log.Print("sender:", sender)

	// Message
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", cfg.FromNick, cfg.FromUser)
	e.To = []string{tos}
	e.Subject = subject
	e.Text = []byte(content)
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")

	var tcfg = new(tls.Config)
	tcfg.ServerName = cfg.ServerHost
	tcfg.InsecureSkipVerify = true

	// Send the email to somebody
	hostPort := fmt.Sprintf("%s:%s", cfg.ServerHost, strconv.Itoa(cfg.ServerPort))
	log.Print("host:port :", hostPort)
	err := e.SendWithTLS(
		hostPort,
		smtp.PlainAuth("", cfg.FromUser, cfg.Password, cfg.ServerHost),
		tcfg,
	)
	if err != nil {
		log.Print("send err:", err)
		w.Write([]byte(`err: send mail failed`))
		return
	}
	// send err:504 5.7.4 Unrecognized authentication
	// type [BJXPR01CA015.CHNPR01.prod.partner.outlook.cn]
	w.Write([]byte(`send mail sucessful`))
}

func main() {
	http.HandleFunc("/mail", handleMail)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
