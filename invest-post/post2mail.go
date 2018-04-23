package main

import (
	"flag"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/BurntSushi/toml"

	gomail "gopkg.in/mail.v2"
	// gopkg.in/gomail.v2 for github.com/go-gomail/gomail
	// - gopkg.in/mail.v2 for github.com/go-mail/mail"
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

// parse post file name
func handleMail(w http.ResponseWriter, r *http.Request) {
	// for parse multipart form, not use r.ParseForm()
	r.ParseMultipartForm(1 << 20)

	var content, subject, sender, tos string
	if r.Method == http.MethodGet {
		content = r.FormValue("content")
		subject = r.FormValue("subject")
		sender = r.FormValue("sender")
		tos = r.FormValue("tos")
		log.Print("hit get!")
	}
	if r.Method == http.MethodPost {
		content = r.PostFormValue("content")
		subject = r.PostFormValue("subject")
		sender = r.PostFormValue("sender")
		tos = r.FormValue("tos")
		log.Print("hit post!")
	}
	w.Write([]byte(fmt.Sprintf("content:%s\n", content)))
	w.Write([]byte(fmt.Sprintf("subject:%s\n", subject)))
	w.Write([]byte(fmt.Sprintf("sender :%s\n", sender)))
	w.Write([]byte(fmt.Sprintf("tos    :%s\n", tos)))

	if len(content) == 0 || len(subject) == 0 ||
		len(tos) == 0 || len(sender) == 0 {
		w.Write([]byte(`err: some parameter is blank`))
		return
	}

	// parse files from form data
	var file multipart.File // an inferface
	var fhdr *multipart.FileHeader
	var err error
	var fname = "hello.md"

	file, fhdr, err = r.FormFile(fname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err: parse file %s failed", fname)))
		return
	}
	log.Printf("file name:%+v\n", file)
	log.Printf("file name:%+v\n", fhdr.Filename)
	log.Printf("file size:%+v\n", fhdr.Size)
	log.Printf("header   :%+v\n", fhdr.Header)

	fname = "golang.md"
	file, fhdr, err = r.FormFile(fname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err: parse file %s failed", fname)))
		return
	}
	log.Printf("file name:%+v\n", file)
	log.Printf("file name:%+v\n", fhdr.Filename)
	log.Printf("file size:%+v\n", fhdr.Size)
	log.Printf("header   :%+v\n", fhdr.Header)

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

	// Send the email to Bob, Cora and Dan.
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
