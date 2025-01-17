@echo off
:: Windows下设置CGO
set CGO_ENABLED=1
set CC=gcc

:: 编译动态库
go build -buildmode=c-shared -o libwsclient.dll ./export/java.go

:: 编译 Android/iOS
::gomobile bind -target=android ./export/mobile.go
::gomobile bind -target=ios ./export/mobile.go

::  编译 WASM
::GOOS=js GOARCH=wasm go build -o main.wasm ./export/wasm.go