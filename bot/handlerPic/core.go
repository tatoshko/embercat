package handlerPic

import "log"

const (
    CDN       = "https://pics.useful.team"
    REDIS_KEY = "pic:anime"
    MAX_PICS  = 409
)

func getLogger(scope string) func(s string) {
    return func(s string) {
        log.Printf("Handler Pic: [%s] %s", scope, s)
    }
}
