package handlerWednesday

import "time"

func ItIsWednesdayMyDudes() bool {
    return !(time.Now().Weekday() == time.Wednesday)
}
