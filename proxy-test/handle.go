package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

const (
	parseErr = "err: parse reader failed"
	paramErr = "err: some parameter is blank"
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

// parse post form and file
func handleMail(w http.ResponseWriter, r *http.Request) {
	formMap, fileMap, merr := parseRequest(w, r)
	if merr != nil {
		return
	}
	logger.Info("formMap:", formMap)
	logger.Info("fileMap:", fileMap)

	to, cc := parseToCc(formMap)
	logger.Print("To:", to)
	logger.Print("Cc:", cc)

	return
}
