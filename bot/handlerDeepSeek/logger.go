package handlerDeepSeek

import (
    "fmt"
    "strings"
)

func getLogger(scope string) func(s ...string) {
    return func(s ...string) {
        fmt.Printf("DeepSeek [%s] %s\n", scope, strings.Join(s, ""))
    }
}
