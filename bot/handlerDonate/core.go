package handlerDonate

import (
    "fmt"
    "strings"
)

const COLLECTOR = "tatoshko"

func getLogger(scope string) func(s ...string) {
    return func(s ...string) {
        fmt.Printf("DONATE [%s] %s", s, strings.Join(s, ""))
    }
}
