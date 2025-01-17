@echo off
:: Windows下设置CGO
set CGO_ENABLED=1
set CC=gcc

:: 编译动态库
go build -buildmode=c-shared -o libwsclient.dll ./export/ws_export.go 