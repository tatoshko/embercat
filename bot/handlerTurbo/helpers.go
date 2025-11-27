package handlerTurbo

import (
    "math/rand"
    "time"
)

func GetRandomLiner() Liner {
    rand.Seed(time.Now().UnixNano())
    n := rand.Int63n(TOTAL_LINERS) + 1

    liner, _ := NewLiner(n, 1)

    return liner
}
