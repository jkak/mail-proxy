package mailproxy

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/jkak/mail"
)

const (
	parseErr  = "parse from reader failed"
	paramErr  = "some parameters are missed"
	s2ChanErr = "send to internal channel timeout"
	connErr   = "connect to remote server timeout"
	sendErr   = "send to remote server timeout"
)

func parseGetRequest(w http.ResponseWriter, r *http.Request) (formMapping, fileMapping, error) {
	r.ParseForm()
	formMap := make(formMapping)

	cont := r.FormValue("content")
	if len(cont) >= cfg.ContLen {
		cont = cont[:cfg.ContLen]
	}
	formMap["content"] = cont
	formMap["subject"] = r.FormValue("subject")
	formMap["sender"] = r.FormValue("sender")
	formMap["tos"] = r.FormValue("tos")
	formMap["cc"] = r.FormValue("cc")

	if len(formMap["content"]) == 0 || len(formMap["subject"]) == 0 ||
		len(formMap["tos"]) == 0 || len(formMap["sender"]) == 0 {
		return nil, nil, fmt.Errorf(paramErr)
	}
	return formMap, nil, nil
}

func parseRequest(w http.ResponseWriter, r *http.Request) (formMapping, fileMapping, error) {
	if r.Method == http.MethodGet {
		return parseGetRequest(w, r)
	}
	rd, err := r.MultipartReader()
	if err != nil {
		Logger.Debugf("parseRequest() %s", parseErr)
		return nil, nil, fmt.Errorf(parseErr)
	}
	formMap := make(formMapping)
	fileMap := make(fileMapping)
	buf := make([]byte, cfg.ContLen)

	var capLen, realLen int
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
			if formName == "content" {
				Logger.Debugf("form info of key=%s, value=<skipped>", formName)
			} else {
				Logger.Debugf("form info of key=%s, value=%s", formName, formMap[formName])
			}
		} else {
			var bb = new(bytes.Buffer)
			io.Copy(bb, part)
			fileMap[fileName] = bb
			Logger.Debugf("file name=%s", fileName)
			capLen += bb.Cap()
			realLen += bb.Len()
		}
	}
	// Logger.Debugf("formMap: %v", formMap)
	Logger.Infof("attachment num=%d; buffer real size=%d, cap size=%d",
		len(fileMap), realLen, capLen)

	// ignore "cc" info
	if len(formMap["content"]) == 0 || len(formMap["subject"]) == 0 ||
		len(formMap["tos"]) == 0 || len(formMap["sender"]) == 0 {
		return formMap, fileMap, fmt.Errorf(paramErr)
	}
	buf = nil
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
	Logger.Debugf("To=%s", to)
	Logger.Debugf("Cc=%s", cc)
	return to, cc
}

func genMsg(formMap formMapping, fileMap fileMapping) (*Sender, *mail.Message) {
	To, Cc := parseToCc(formMap)

	// Message
	msg := mail.NewMessage()
	msg.SetHeader("To", To...)
	if len(Cc) != 0 && Cc[0] != "" {
		msg.SetHeader("Cc", Cc...)
	}
	msg.SetHeader("Subject", formMap["subject"])
	msg.SetBody("text/html", formMap["content"])
	// attachments
	for fName, contReader := range fileMap {
		msg.AttachReader(fName, contReader)
	}

	var sender Sender
	senders := senderMap[formMap["sender"]]
	var idx int32
	if len(senders) != 1 {
		idx = rand.Int31n(int32(len(senders)))
	}
	sender = senders[idx]
	msg.SetHeader("From", sender.MailUser)
	Logger.Debugf("select sender: rand=%d; sender=%+v", idx, sender.MailUser)
	return &sender, msg
}

func sendMail(formMap formMapping, fileMap fileMapping) error {
	sender, msg := genMsg(formMap, fileMap)

	msgPkt := &sendMsg{
		msg:      msg,
		statusCh: make(chan string),
		mailUser: sender.MailUser,
	}

	// write to chan with timeout check
	select {
	case <-time.After(time.Second * time.Duration(cfg.Timeout)):
		return fmt.Errorf(s2ChanErr)
	case mailCh[sender.MailUser] <- msgPkt:
		Logger.Debug("sendMail() send msg to mail channel success")
	}
	// wait goroutine send msg finished
	select {
	case <-time.After(time.Second * time.Duration(cfg.Timeout)):
		return fmt.Errorf(sendErr)
	case <-msgPkt.statusCh:
		return nil
	}
	return fmt.Errorf(sendErr)
}

// parse post form and file
func handleMail(w http.ResponseWriter, r *http.Request) {
	resp := new(resp2User)
	resp.Status = cfg.StatusOK

	formMap, fileMap, merr := parseRequest(w, r)
	if merr != nil {
		resp.Msg = merr.Error()
		resp.Status = cfg.StatusEr
		goto RESPONSE
	}
	if err := sendMail(formMap, fileMap); err != nil {
		resp.Status = cfg.StatusEr
		resp.Msg = err.Error()
	}

RESPONSE:
	resp.Time = time.Now().Unix()
	writeResp(w, resp)
}
