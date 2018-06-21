package main

import (
	"bytes"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	gomail "github.com/jkak/mail"
)

const (
	parseErr = "err: parse reader failed"
	paramErr = "err: some parameter is blank"
	statusEr = "err: send mail failed"
	statusOK = "send mail sucessful"
)

type formMapping map[string]string
type fileMapping map[string]*bytes.Buffer

func parseRequest(w http.ResponseWriter, r *http.Request) (formMapping, fileMapping, error) {
	formMap := make(formMapping)
	fileMap := make(fileMapping)
	buf := make([]byte, cfg.ContLen)

	rd, err := r.MultipartReader()
	if err != nil {
		return formMap, fileMap, errors.New(parseErr)
	}
	for {
		part, perr := rd.NextPart()
		if perr == io.EOF {
			break
		}
		fileName := part.FileName()
		formName := part.FormName()

		if fileName == "" {
			n, _ := part.Read(buf)
			formMap[formName] = string(buf[0:n])
		} else {
			var bb bytes.Buffer
			fileMap[fileName] = &bb
			io.Copy(fileMap[fileName], part)
		}
	}

	// ignore "cc" info
	if len(formMap["content"]) == 0 || len(formMap["subject"]) == 0 ||
		len(formMap["tos"]) == 0 || len(formMap["sender"]) == 0 {
		w.Write([]byte(paramErr))
		return formMap, fileMap, errors.New(paramErr)
	}
	return formMap, fileMap, nil
}

func parseToCc(f formMapping) ([]string, []string) {
	// to and cc
	var to, cc []string
	for _, t := range strings.Split(f["tos"], ",") {
		to = append(to, strings.TrimSpace(t))
	}
	for _, c := range strings.Split(f["cc"], ",") {
		cc = append(cc, strings.TrimSpace(c))
	}
	return to, cc
}

func sendMail(formMap formMapping, fileMap fileMapping) error {
	To, Cc := parseToCc(formMap)
	logger.Print("To:", To)
	logger.Print("Cc:", Cc)

	// Message
	m := gomail.NewMessage()
	m.SetHeader("To", To...)
	if len(Cc) != 0 && Cc[0] != "" {
		m.SetHeader("Cc", Cc...)
	}
	m.SetHeader("Subject", formMap["subject"])
	m.SetBody("text/html", formMap["content"])
	// attachments
	for fname, contReader := range fileMap {
		m.AttachReader(fname, contReader)
	}

	var sender Sender
	senders := senderMap[formMap["sender"]]
	var idx int32
	if len(senders) != 1 {
		idx = rand.Int31n(int32(len(senders)))
	}
	sender = senders[idx]
	m.SetHeader("From", sender.Account)
	logger.Printf("rand:%d; sender=%+v", idx, sender.Account)

	d := gomail.NewDialer(cfg.ServerHost, cfg.ServerPort, sender.Account, sender.Password)
	for i := 1; i <= 3; i++ {
		err := d.DialAndSend(m)
		if err != nil {
			logger.Printf("send err, retry=%d, err info:%s\n", i, err)
			time.Sleep(time.Second * time.Duration(cfg.ProxySleep))
			continue
		} else {
			logger.Print(statusOK)
			return nil
		}
	}
	return errors.New(statusEr)
}

// parse post form and file
func handleMail(w http.ResponseWriter, r *http.Request) {
	formMap, fileMap, merr := parseRequest(w, r)
	if merr != nil {
		return
	}
	logger.Info("formMap:", formMap)
	logger.Info("fileMap:", fileMap)

	if err := sendMail(formMap, fileMap); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(statusOK))
	}
}
