package redis

import (
    "gopkg.in/redis.v3"
)

func GetClient() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr: "127.0.0.1",
        DB:   1,
    })
}
