## proxy-test

server

```shell
# go run *go -c ./proxy.toml
sender: alarm	Len:3
sender: monitor	Len:1
sender: test	Len:1

```



client

```shell
# python sendPost.py op@op.cn "" test
send mail sucessful
200

# python sendHtml.py op@op.cn "" test
send mail sucessful
200
```

