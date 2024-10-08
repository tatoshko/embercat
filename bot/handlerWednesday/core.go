package handlerWednesday

import (
    "fmt"
    "log"
    "strings"
)

const (
    NO_WEDNESDAY = "no-wednesday.jpg"
)

func getLogger(scope string) func(s ...string) {
    return func(s ...string) {
        log.Printf(fmt.Sprintf("WEDNESDAY [%s] %s", scope, strings.Join(s, "")))
    }
}
