# IpParser

ip\_parser是一款Go实现的纯真IP库解析库

[![License](https://img.shields.io/:license-apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/ant-libs-go/ip_parser?status.png)](http://godoc.org/github.com/ant-libs-go/ip_parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/ant-libs-go/ip_parser)](https://goreportcard.com/report/github.com/ant-libs-go/ip_parser)

## 特性

* 支持将纯真库加载到内存中进行查找，提高效率
* 对纯真库返回的非标准格式进行尽可能的标准化

## 安装

	go get github.com/ant-libs-go/ip_parser

## 快速开始

```golang
dat, _ := ioutil.ReadFile("./qqwry.dat")

fmt.Printf("%+v\n", NewIpParser("182.99.99.91", dat).Parse())
```
