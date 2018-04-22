

## 测试test

### 1) go-gomail/gomail 

github.com/go-gomail/gomail

```shell
# conf: 
#	test.toml for bublic
# 	mail.toml for myself
go run go-gomail-gomail.go -c mail.toml

```

curl

```shell
curl -v  "localhost:8080/mail?content=my-test&subject=hello&tos=jkak@163.com&sender=alarm"

< HTTP/1.1 200 OK
< Date: Sun, 22 Apr 2018 14:32:18 GMT
< Content-Length: 19
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
send mail sucessful%
```



### 2) jordan-wright/email

github.com/jordan-wright/email

```shell
go run jordan-wright-email.go -c mail.toml

# 2018/04/22 23:22:32 host:port :smtp.partner.outlook.cn:587
# 2018/04/22 23:22:42 send err:504 5.7.4 Unrecognized authentication type [BJXPR01CA015.CHNPR01.prod.partner.outlook.cn]

```

curl

```shell
curl -v  "localhost:8080/mail?content=my-test&subject=hello&tos=jkak@163.com&sender=alarm"

< HTTP/1.1 200 OK
< Date: Sun, 22 Apr 2018 14:32:18 GMT
< Content-Length: 19
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
err: send mail failed%
```

虽然使用了TLS老地方，但依然会报错，且报错信息同使用telnet报错相同。

报错信息：`504 5.7.4 Unrecognized authentication type`





### 3) domodwyer/mailyak

github.com/domodwyer/mailyak



```shell
go run domodwyer-mailyak.go -c mail.toml

# 2018/04/22 23:41:52 host:port :smtp.partner.outlook.cn:587
# 2018/04/22 23:42:03 send err:504 5.7.4 Unrecognized authentication type [BJXPR01CA015.CHNPR01.prod.partner.outlook.cn]
```

curl

```shell
curl -v  "localhost:8080/mail?content=my-test&subject=hello&tos=jkak@163.com&sender=alarm"

< HTTP/1.1 200 OK
< Date: Sun, 22 Apr 2018 14:32:18 GMT
< Content-Length: 19
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
err: send mail failed%
```

10秒超时后，也会出现同上一个库一样的报错信息：`504 5.7.4 Unrecognized authentication type`

