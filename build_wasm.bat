@echo off
:: 设置 GOROOT 环境变量
set GOROOT=E:\develop\go

:: 设置 WASM 环境变量
set GOOS=js
set GOARCH=wasm

:: 设置输出目录
set WASM_DIR=D:\code\github\go-ws-client\wasm

:: 创建 wasm 目录（如果不存在）
if not exist "%WASM_DIR%" mkdir "%WASM_DIR%"

:: 编译 WASM
go build -tags wasm -o "%WASM_DIR%\main.wasm" ./export/wasm.go

:: 复制 wasm_exec.js
copy "%GOROOT%\misc\wasm\wasm_exec.js" "%WASM_DIR%\"

echo.
echo WASM files built successfully in %WASM_DIR%
