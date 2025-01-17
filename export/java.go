package main

/*
#include <stdlib.h>
typedef void (*MessageCallback)(char*);

static void invokeCallback(MessageCallback callback, char* message) {
    if (callback != NULL) {
        callback(message);
    }
}
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/imai/go-ws-client/client"
)

var (
	globalClient    *client.WSClient
	messageCallback C.MessageCallback
)

//export InitWebSocket
func InitWebSocket(url, token, deviceType *C.char) C.int {
	urlStr := C.GoString(url)
	tokenStr := C.GoString(token)
	deviceTypeStr := C.GoString(deviceType)

	globalClient = client.NewWSClient(urlStr, tokenStr, deviceTypeStr)

	// 如果有回调，设置它
	if messageCallback != nil {
		globalClient.SetMessageCallback(func(message string) {
			cMessage := C.CString(message)
			C.invokeCallback(messageCallback, cMessage)
			C.free(unsafe.Pointer(cMessage))
		})
	}

	err := globalClient.Connect()
	fmt.Println("InitWebSocket", err)
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

//export SetMessageCallback
func SetMessageCallback(callback C.MessageCallback) {
	messageCallback = callback
	if globalClient != nil {
		globalClient.SetMessageCallback(func(message string) {
			cMessage := C.CString(message)
			C.invokeCallback(messageCallback, cMessage)
			C.free(unsafe.Pointer(cMessage))
		})
	}
}

func main() {}
