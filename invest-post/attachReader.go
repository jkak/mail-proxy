package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"

	gomail "github.com/jkak/mail"
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

// parse post form and file
func handleMail(w http.ResponseWriter, r *http.Request) {
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

	formMap := make(map[string]string)
	fileMap := make(map[string]*bytes.Buffer)
	buf := make([]byte, 1024*1024)

	for {
		// Reader is an iterator over parts in a MIME multipart body.
		// Reader's underlying parser consumes its input as needed.
		// NextPart() returns the next part in the multipart or an error.
		part, perr := rd.NextPart()
		// When there are no more parts, the error io.EOF is returned
		if perr == io.EOF {
			break
		}

		fileName := part.FileName()
		formName := part.FormName()

		log.Printf("file header:%+v\n", part.Header)
		if fileName == "" {
			n, _ := part.Read(buf)
			formMap[formName] = string(buf[0:n])
			log.Printf("form name:%s\n", formName)
			log.Printf("form cont:%s\n", formMap[formName])
		} else {
			var bb bytes.Buffer
			fileMap[fileName] = &bb
			io.Copy(fileMap[fileName], part)
			log.Printf("file name:%s\n", fileName)
			log.Printf("file cont:%s\n", fileMap[fileName].String())
		}
	}
	// for test here
	//for k, v := range fileMap {
	//	w.Write([]byte(k))
	//	io.Copy(w, v)
	//}

	if len(formMap["content"]) == 0 || len(formMap["subject"]) == 0 ||
		len(formMap["tos"]) == 0 || len(formMap["sender"]) == 0 {
		w.Write([]byte(`err: some parameter is blank`))
		return
	}

	//return

	// Message
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.FromUser)

	m.SetHeader("To", formMap["tos"])
	if formMap["cc"] != "" {
		m.SetHeader("Cc", formMap["cc"])
	}
	m.SetHeader("Subject", formMap["subject"])
	m.SetBody("text/html", formMap["content"])
	log.Print("sender:", formMap["sender"])

	for fname, contReader := range fileMap {
		m.AttachReader(fname, contReader)
	}

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
