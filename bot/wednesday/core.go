package wednesday

import (
    "fmt"
    "log"
)

const (
    REDIS_KEY    = "frogs"
    NO_WEDNESDAY = "no-wednesday.jpg"
)

func getLogger(scope string) func(s string) {
    return func(s string) {
        log.Printf(fmt.Sprintf("WEDNESDAY [%s] %s", scope, s))
    }
}
