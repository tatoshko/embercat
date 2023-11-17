package handlerTurbo

import (
    "errors"
    "fmt"
    "gopkg.in/redis.v3"
    "log"
    "time"
)

type Todayer struct {
    connection  *redis.Client
    redisKey    string
    currentDate string
}

func NewTodayer(redis *redis.Client, userID int64) (todayer Todayer) {
    return Todayer{
        connection:  redis,
        redisKey:    makeUserCollectionKey(REDIS_KEY_TURBO_DAY, userID),
        currentDate: time.Now().Format("2006-Jan-02"),
    }
}

func (t Todayer) Dirty() (yes bool, err error) {
    if yes, err = t.connection.SIsMember(t.redisKey, t.currentDate).Result(); err != nil {
        log.Printf("HandlerTurbo SIsMember error %s", err.Error())
        return false, errors.New(fmt.Sprintf("can't check todays activity for %s", t.redisKey))
    }

    if yes {
        return true, nil
    }

    if _, err = t.connection.SAdd(t.redisKey, t.currentDate).Result(); err != nil {
        log.Printf("HandlerTurbo SAdd error %s", err.Error())
        return false, errors.New(fmt.Sprintf("can't set todays activity for %s", t.redisKey))
    }

    return false, nil
}
