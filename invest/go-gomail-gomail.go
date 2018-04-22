// https://github.com/go-gomail/gomail/
// https://godoc.org/gopkg.in/gomail.v2#example-package
// https://godoc.org/github.com/go-gomail/gomail#example-package--Daemon
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/go-gomail/gomail"
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
	// Message
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.FromUser)

	m.SetHeader("To", tos)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)
	//m.Attach("~/Desktop/saas-iaas-paas.jpeg")
	//m.Attach("~/Desktop/open-falcon.jpeg")
	log.Print("sender:", sender)
	d := gomail.NewDialer(cfg.ServerHost, cfg.ServerPort, cfg.FromUser, cfg.Password)

	// Send the email to somebody
	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
		w.Write([]byte(`err: send mail failed`))
		return
	}
	w.Write([]byte(`send mail sucessful`))
}

func main() {
	http.HandleFunc("/mail", handleMail)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
