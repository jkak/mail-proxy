#!/usr/bin/env python
# _*_ coding: utf-8 _*_

import requests
import os, sys

if len(sys.argv) != 4:
    print len(sys.argv)
    print "usage: %s tos cc sender" % sys.argv[0]
    sys.exit(1)

tos = sys.argv[1]
cc = sys.argv[2]
sender = sys.argv[3]

file_names = ['text/china中国.txt', 'text/golang.md', 'text/hello.md']
base = os.path.basename
# (formname, (filename, fileDATA))
file_data = [(base(name), (name, open(name, 'rb').read())) for name in file_names]

url = 'http://localhost:8080/mail'

with open('text/index.html') as fd:
    msg = fd.read()

data = {
   "subject": "subject-here",
   'tos': tos,
   'cc': cc,
   'mailtype': 'html',
   'sender': sender,
   'content': msg
}

r = requests.post(url, files=file_data, data=data)
print(r.text)
print(r.status_code)
