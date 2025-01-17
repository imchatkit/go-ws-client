//go:build wasm

package main

import "syscall/js"

// RegisterWASMFunctions 注册WASM函数
func RegisterWASMFunctions() {
	client := &WSClientWrapper{}

	js.Global().Set("wsConnect", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return false
		}
		return client.Connect(args[0].String(), args[1].String(), args[2].String())
	}))

	// ... 其他函数注册
}
