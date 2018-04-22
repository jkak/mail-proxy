package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/domodwyer/mailyak"
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
	tos := r.FormValue("tos")
	sender := r.FormValue("sender")

	if len(content) == 0 || len(subject) == 0 ||
		len(tos) == 0 || len(sender) == 0 {
		w.Write([]byte(`err: some paramter is blank`))
		return
	}
	hostPort := fmt.Sprintf("%s:%s", cfg.ServerHost, strconv.Itoa(cfg.ServerPort))
	log.Print("host:port :", hostPort)
	// Message
	m := mailyak.New(hostPort,
		smtp.PlainAuth("", cfg.FromUser, cfg.Password, cfg.ServerHost),
	)

	m.To(tos)
	m.From(cfg.FromUser)
	m.FromName(cfg.FromNick)
	m.Subject(subject)
	m.Plain().Set(content)
	//m.HTML().Set(content)

	// sth wrong when use Attach() method
	//fd, _ := os.Open("~/Desktop/saas-iaas-paas.jpeg")
	//m.Attach("saas-iaas-paas.jpeg", fd)
	log.Print("sender:", sender)

	// Send the email to Bob, Cora and Dan.
	if err := m.Send(); err != nil {
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
