package utils

import (
	"log"
	"runtime/debug"
)

func CatchPanic(responder *interface{}) {
	if r := recover(); r != nil {
		log.Printf("panic recovered: %v\n%s", r, debug.Stack())
		if responder != nil {
			*responder = nil
		}
	}
}

