# Migration

一个简单的图片迁移工具。

因为博客以前使用的图床不稳定，博客经常加载不出图片，于是打算把博客里的图片换一个图床，但是自己去手动下载上传又太过麻烦，于是写了哥小工具来批量处理。

具体做的事情就是 读取 Markdown 文件，使用正则匹配出图片的链接，筛选出需要转移的图片，下载并上传到指定的其他图床。并且把 Markdown 中对应的链接改成新图片的链接。

适用场景

- md 中图片批量转移
- md 中图片批量下载到本地

## Setup

```bash
$ cd little-tools/migration
$ go build -o mg .
$ ./mg --help

  -cfg string
   	[optional] config file, if this specified, other flags won't be used.
 -down string
   	[*] the base url that you download picture. ex: jianshu.io
 -fk string
   	[optional] fileKey, use default if you dont know what it mean. (default "file")
 -if string
   	[*] input file
 -of string
   	[optional] output file
 -re string
   	[optional] regular expression used to match link. (default "http[s?]://[\\d\\w\\./-]*\\.(jpg|png)")
 -tk string
   	[optional] tokenKey, use default if you dont know what it mean. (default "token")
 -token string
   	[optional] the token if server needed (default "token")
 -up string
   	[*] the url that upload picture. ex: http://localhost:8000/upload
```

## Examples

### Single File

将 xxx.md 中的图片下载后上传到 http://xxx.xxxx.xxx/upload ，替换文件中的链接后输出到 xxx_new.md 中。<br />你也可以自己指定输出文件的名字。

```bash
./mg -if=./xxx.md -up="http://xxx.xxxx.xxx/upload" -down="jianshu.io" -token=password
```

### All .md Files In Directory

如果输入文件指定为一个文件夹，则会读取文件夹中所有 .md 文件并且处理生成对应的 *_new.md 文件。

```bash
./mg -if=./posts -up="http://localhost:8000/upload" -down="jianshu.io" -token=password
```
