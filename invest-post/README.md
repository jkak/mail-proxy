



### 1 about gomail

在上述三个库中。gomail是star上千，且使用也比较方便的库。并且有示例可以做为daemon运行。但这个库最大的问题是作者已经有2年没有更新了。而且从issues中的信息看。作者的邮箱也已经满了而没有处理。感觉作者已经彻底消失了。

已经有两三个用户从gomail开分支进行后继维护。但从目前的进度来看，新增加的功能还不够多。也没有较好的增加对附件文件的流式处理。而这个功能，是将gomail用做本地代理时强依赖的功能点。

原库：https://github.com/go-gomail/gomail

新的分支：https://github.com/go-mail/mail



### 2 test mail for post

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




### 3 test reader for post

server

```shell
# go run multipart-reader.go -c mail.toml

2018/04/24 23:41:27 file name:
2018/04/24 23:41:27 form name:content
2018/04/24 23:41:27 	file header:map[Content-Disposition:[form-data; name="content"]]
2018/04/24 23:41:27 	conent len:28
2018/04/24 23:41:27 	conent buf:content here without newline
2018/04/24 23:41:27 file name:
2018/04/24 23:41:27 form name:tos
2018/04/24 23:41:27 	file header:map[Content-Disposition:[form-data; name="tos"]]
2018/04/24 23:41:27 	conent len:8
2018/04/24 23:41:27 	conent buf:op@jk.cn
2018/04/24 23:41:27 file name:
2018/04/24 23:41:27 form name:sender
2018/04/24 23:41:27 	file header:map[Content-Disposition:[form-data; name="sender"]]
2018/04/24 23:41:27 	conent len:2
2018/04/24 23:41:27 	conent buf:op
2018/04/24 23:41:27 file name:
2018/04/24 23:41:27 form name:mailtype
2018/04/24 23:41:27 	file header:map[Content-Disposition:[form-data; name="mailtype"]]
2018/04/24 23:41:27 	conent len:4
2018/04/24 23:41:27 	conent buf:text
2018/04/24 23:41:27 file name:
2018/04/24 23:41:27 form name:subject
2018/04/24 23:41:27 	file header:map[Content-Disposition:[form-data; name="subject"]]
2018/04/24 23:41:27 	conent len:12
2018/04/24 23:41:27 	conent buf:subject-here
2018/04/24 23:41:27 file name:text/golang.md
2018/04/24 23:41:27 form name:golang.md
2018/04/24 23:41:27 	file header:map[Content-Disposition:[form-data; name="golang.md"; filename="text/golang.md"]]
2018/04/24 23:41:27 	conent len:13
2018/04/24 23:41:27 	conent buf:hello golang

2018/04/24 23:41:27 file name:text/hello.md
2018/04/24 23:41:27 form name:hello.md
2018/04/24 23:41:27 	file header:map[Content-Disposition:[form-data; name="hello.md"; filename="text/hello.md"]]
2018/04/24 23:41:27 	conent len:32
2018/04/24 23:41:27 	conent buf:## test markdown
hello markdown
```

client

```shell
# python sendPost.py op@jk.cn op

200
```

如上，通过r.MultipartReader()，可以获取到一个reader，并通过遍历reader，可以得到所有的post信息。包括：

* data相关字段
  * 对应python中request.post中的data参数。
  * 这些字段只有form name，没有file name。
  * 通过p.FormName()可以获得相关字段：
    * subject
    * content
    * tos
    * sender
  * 通过对应的p.Read(buf)可以将内容读到buf中。
* files相关字段
  * 对应python中request.post中的files参数。
  * files字段即有form name，也有file name。
    * form name对应请求来源中file_data列表中的第一个名字；
    * file name对应请求来源中file_data列表中的第二个名字；
  * 通过p.FileName()可以获得相关文件名：
    * text/golang.md
    * text/hello.md
  * 对应的p.Read(buf)可以将内容读到buf中。



### 4 parse post from reader

server

```shell
# go run parseReader.go -c mail.toml
2018/04/25 09:37:28 reader:&{bufReader:0xc4200b62a0 currentPart:<nil> partsRead:0 nl:[13 10] nlDashBoundary:[13 10 45 45 97 102 48 52 54 102 100 99 53 49 48 51 52 50 50 98 98 99 53 102 102 49 50 52 48 55 97 100 101 52 50 53] dashBoundaryDash:[45 45 97 102 48 52 54 102 100 99 53 49 48 51 52 50 50 98 98 99 53 102 102 49 50 52 48 55 97 100 101 52 50 53 45 45] dashBoundary:[45 45 97 102 48 52 54 102 100 99 53 49 48 51 52 50 50 98 98 99 53 102 102 49 50 52 48 55 97 100 101 52 50 53]}
2018/04/25 08:37:28 form name:content; buf:content here without newline
2018/04/25 08:37:28 form name:tos; buf:op@jk.cn
2018/04/25 08:37:28 form name:sender; buf:op
2018/04/25 08:37:28 form name:mailtype; buf:text
2018/04/25 08:37:28 form name:subject; buf:subject-here
2018/04/25 08:37:28 	file name: text/golang.md
2018/04/25 08:37:28 	file name: text/hello.md
k=mailtype, v=text
k=subject, v=subject-here
k=content, v=content here without newline
k=tos, v=op@jk.cn
k=sender, v=op
```

client

```shell
# python sendPost.py op@jk.cn op
mailtype:text
subject:subject-here
content:content here without newline
tos:op@jk.cn
sender:op
file num:2

200
```

如上，通过对reader进行NextPart()循环处理。可以获取到所有的请求参数及文件。

将form参数放入formMap中，file参数放入fileMap中。

剩下的问题是，能否将接收到的文件流，直接发送出去，而不是将其先存盘后再调用Attach()方法再次读取。



