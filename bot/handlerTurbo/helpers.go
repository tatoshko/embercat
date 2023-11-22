package handlerTurbo

import (
    "math/rand"
    "time"
)

func GetRandomLiner() Liner {
    rand.Seed(time.Now().UnixNano())
    liner, _ := NewLiner(rand.Int63n(TOTAL_LINERS), 1)
    return liner
}
