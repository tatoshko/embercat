package handlerPic

import "log"

func getLogger(scope string) func(s string) {
    return func(s string) {
        log.Printf("Handler Pic: [%s] %s", scope, s)
    }
}
