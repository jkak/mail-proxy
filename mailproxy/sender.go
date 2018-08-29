package mailproxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jkak/mail"
)

type sendMsg struct {
	msg      *mail.Message
	statusCh chan string
	mailUser string
}

type resp2User struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Time   int64  `json:"time"`
}

func writeResp(w http.ResponseWriter, resp interface{}) {
	buf, err := json.Marshal(resp)
	if err != nil {
		Logger.Errorf("write response error: %v", err.Error())
	}
	w.Write(buf)
}

func mailSenders() {
	for _, sd := range cfg.Senders {
		go send2Server(sd)
	}
}

// send msg to smtp server
func send2Server(sender Sender) {
	var sc mail.SendCloser
	var err error
	open := false

	for {
	OUT_FOR:
		select {
		case msg, ok := <-mailCh[sender.MailUser]:
			if !ok {
				return // channel closed
			}
			if !open {
				sc, err = connect(sender)
				if err != nil {
					sc, _ = reConnect(sender)
				}
				open = true
			}
			// retry
			for i := 1; i <= cfg.RetryTimes; i++ {
				if err = mail.Send(sc, msg.msg); err != nil {
					Logger.Printf("send err, retry=%d, err info:%s\n", i, err)
					time.Sleep(time.Second * time.Duration(cfg.RetryInterval))
					continue
				} else {
					Logger.Printf("send2Server() send mail result: %s", cfg.StatusOK)
					msg.statusCh <- cfg.StatusOK
					goto OUT_FOR
				}
			}
			msg.statusCh <- cfg.StatusEr

		case <-time.After(time.Second * time.Duration(cfg.ProxyKeepalive)):
			if open {
				sc.Close()
				open = false
			}
		}
	}
}

func connect(s Sender) (mail.SendCloser, error) {
	dstServer := mail.NewDialer(cfg.ServerHost, cfg.ServerPort, s.MailUser, s.Password)
	baseGap := 500 * time.Millisecond
	for {
		s, err := dstServer.Dial()
		if err != nil {
			Logger.Error("dial to remote server failed")
			time.Sleep(baseGap)
			baseGap *= 2
			if baseGap > time.Second*time.Duration(cfg.Timeout) {
				baseGap = time.Second * time.Duration(cfg.Timeout)
			}
			continue
		}
		return s, nil
	}
	return nil, fmt.Errorf(connErr)
}

func reConnect(s Sender) (mail.SendCloser, error) {
	sc, err := connect(s)
	return sc, err
}
