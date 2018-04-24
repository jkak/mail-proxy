package main

import (
	"flag"
	"io"
	"log"
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
	// parse reader from form data
	//var rd *multipart.Reader
	//var fhdr *multipart.FileHeader

	// func (r *Request) MultipartReader() (*multipart.Reader, error)
	// MultipartReader returns a MIME multipart reader if this is a
	// multipart/form-data POST request, else returns nil and an error.
	// Use this function instead of ParseMultipartForm
	// to process the request body as a stream.
	rd, err := r.MultipartReader()
	if err != nil {
		w.Write([]byte(`err: parse reader failed`))
		log.Print("err: parse reader failed")
		return
	}
	log.Printf("reader:%+v\n", rd)

	var content, subject, sender, tos string
	for {
		// Reader is an iterator over parts in a MIME multipart body.
		// Reader's underlying parser consumes its input as needed.
		// NextPart() returns the next part in the multipart or an error.
		part, perr := rd.NextPart()
		// When there are no more parts, the error io.EOF is returned
		if perr == io.EOF {
			break
		}

		// filename param of the Part's C-D header.
		log.Printf("file name:%s\n", part.FileName())
		// name param if p has a C-D of type "form-data". or empty string
		log.Printf("form name:%s\n", part.FormName())

		//if part.FileName() != "" {
		log.Printf("\tfile header:%+v\n", part.Header)
		buf := make([]byte, 1024)
		n, _ := part.Read(buf)
		log.Printf("\tconent len:%d\n", n)
		log.Printf("\tconent buf:%s\n", string(buf))
		//}
	}

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
