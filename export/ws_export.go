package main

/*
#include <stdlib.h>
typedef int int32_t;
*/
import "C"
import (
	"github.com/imai/go-ws-client/client"
)

var globalClient *client.WSClient

//export InitWebSocket
func InitWebSocket(url, token, deviceType *C.char) C.int {
	urlStr := C.GoString(url)
	tokenStr := C.GoString(token)
	deviceTypeStr := C.GoString(deviceType)

	globalClient = client.NewWSClient(urlStr, tokenStr, deviceTypeStr)
	err := globalClient.Connect()
	if err != nil {
		return C.int(0)
	}
	return C.int(1)
}

//export SendMessage
func SendMessage(message *C.char) C.int {
	if globalClient == nil {
		return C.int(0)
	}

	messageStr := C.GoString(message)
	err := globalClient.Send(messageStr)
	if err != nil {
		return C.int(0)
	}
	return C.int(1)
}

//export CloseWebSocket
func CloseWebSocket() C.int {
	if globalClient == nil {
		return C.int(1)
	}

	err := globalClient.Close()
	globalClient = nil
	if err != nil {
		return C.int(0)
	}
	return C.int(1)
}

func main() {}
