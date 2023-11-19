package handlerPic

import (
    redis2 "embercat/redis"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func RndShare(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("LOAD")

    redis := redis2.GetClient()
    if redis == nil {
        return
    }
    defer redis.Close()

    var pic string
    if pic, err = redis.SRandMember(REDIS_KEY).Result(); err != nil {
        logger(err.Error())
    }

    msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, pic)
    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
