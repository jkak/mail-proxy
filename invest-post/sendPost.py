#!/usr/bin/env python
# _*_ coding: utf-8 _*_

import requests
import os, sys

if len(sys.argv) != 3:
    print len(sys.argv)
    print "usage: %s tos sender" % sys.argv[0]
    sys.exit(1)

tos = sys.argv[1]
sender = sys.argv[2]

file_names = ['text/golang.md', 'text/hello.md']
base = os.path.basename
# (formname, (filename, fileDATA))
file_data = [(base(name), (name, open(name, 'rb').read())) for name in file_names]

url = 'http://localhost:8080/mail'

with open('text/cont.txt') as fd:
    msg = fd.read()

data = {
   "subject": "subject-here",
   'tos': tos,
   'mailtype': 'text',
   'sender': sender,
   'content': msg
}

r = requests.post(url, files=file_data, data=data)
print(r.text)
print(r.status_code)
