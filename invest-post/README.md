



### about gomail

在上述三个库中。gomail是star上千，且使用也比较方便的库。并且有示例可以做为daemon运行。但这个库最大的问题是作者已经有2年没有更新了。而且从issues中的信息看。作者的邮箱也已经满了而没有处理。感觉作者已经彻底消失了。

已经有两三个用户从gomail开分支进行后继维护。但从目前的进度来看，新增加的功能还不够多。也没有较好的增加对附件文件的流式处理。而这个功能，是将gomail用做本地代理时强依赖的功能点。

原库：https://github.com/go-gomail/gomail

新的分支：https://github.com/go-mail/mail



### test mail for post

https://github.com/go-mail/mail forked form gomail.

server:

```shell
# go run post2mail.go -c mail.toml
2018/04/23 23:17:53 cfg:{ServerHost:smtp.partner.outlook.cn ServerPort:587 FromUser:op.cn FromNick:op Password:op}

2018/04/23 23:38:01 hit post!
2018/04/23 23:38:01 file name:{SectionReader:0xc42008da40}
2018/04/23 23:38:01 file name:text/hello.md
2018/04/23 23:38:01 file size:32
2018/04/23 23:38:01 header   :map[Content-Disposition:[form-data; name="hello.md"; filename="text/hello.md"]]
2018/04/23 23:38:01 file name:{SectionReader:0xc42008daa0}
2018/04/23 23:38:01 file name:text/golang.md
2018/04/23 23:38:01 file size:13
2018/04/23 23:38:01 header   :map[Content-Disposition:[form-data; name="golang.md"; filename="text/golang.md"]]
2018/04/23 23:38:01 sender:op

```

client

```shell
# python sendPost.py jk@op.cn op
src_golang.md src_golang.md
src_hello.md src_hello.md
content:content here without newline
subject:subject-here
sender :op
tos    :jk@op.cn
send mail sucessful
200

# wc -c text/*
  28 text/cont.txt
  13 text/golang.md
  32 text/hello.md
  73 total
```

上述post请求的结果，有几个方面的效果：

* post上传的文件已经可以解析到文件名。即使带目录的情况下，也能正常解析
* 后端通过`*multipart.FileHeader`结构体可以获得文件的相关信息
  * 文件名：fhdr.Filename
  * 文件大小：fhdr.Size
  * 文件头：fhdr.Header。可以获取上python post时的简写名称及完整名称。
* 后端收到的文件大小，符合文件本身的属性。
* 最重要的，邮件能够发送出来。

问题：

* go-mail/mail包做为客户端使用，是可以直接带附件，但做为代理端，因为通过sendPost.py上传的文件，已经变成了文件流，因此不能直接使用Attach方法直接发送。

