package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/BurntSushi/toml"

	gomail "gopkg.in/mail.v2"
	// github.com/go-mail/mail
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
	// func (r *Request) MultipartReader() (*multipart.Reader, error)
	rd, err := r.MultipartReader()
	if err != nil {
		w.Write([]byte(`err: parse reader failed`))
		log.Print("err: parse reader failed")
		return
	}
	log.Printf("reader:%+v\n", rd)

	var content, subject, sender, tos string
	buf := make([]byte, 4096)
	fileMap := make(map[string]*multipart.Part)
	formMap := make(map[string]string)

	for {
		// Reader's underlying parser consumes its input as needed.
		part, perr := rd.NextPart()
		if perr == io.EOF {
			break
		}

		if part.FileName() == "" {
			// process form info
			n, _ := part.Read(buf)
			formMap[part.FormName()] = string(buf[0:n])
			fn := part.FormName()
			partStr := string(buf[0:n])
			log.Printf("form name:%s; buf:%s\n", fn, partStr)
		} else {
			log.Printf("\tfile name: %s\n", part.FileName())
			fileMap[part.FileName()] = part
		}
	}
	for k, v := range formMap {
		fmt.Printf("k=%+v, v=%+v\n", k, v)
		w.Write([]byte(fmt.Sprintf("%s:%s\n", k, v)))
	}

	if len(formMap["content"]) == 0 || len(formMap["subject"]) == 0 ||
		len(formMap["tos"]) == 0 || len(formMap["sender"]) == 0 {
		w.Write([]byte(`err: some parameter is blank`))
		return
	}
	w.Write([]byte(fmt.Sprintf("file num:%d\n", len(fileMap))))
	return

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
