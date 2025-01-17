package main

/*
#include <stdint.h>
*/
import "C"

//export WSConnect
func WSConnect(url, token, deviceType *C.char) C.int32_t {
	client := NewWSClient(C.GoString(url), C.GoString(token), C.GoString(deviceType))
	return C.int32_t(boolToInt(client.Connect()))
}

// ... 其他导出函数
