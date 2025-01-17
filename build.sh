#!/bin/bash

# Windows下设置CGO
export CGO_ENABLED=1
export CC=gcc

# 编译动态库
go build -buildmode=c-shared -o libwsclient.dll ./export/ws_export.go 
#用这个也可以 go build -buildmode=c-shared -o libwsclient.so ./client/... 