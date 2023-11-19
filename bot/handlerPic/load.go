package handlerPic

import (
    redis2 "embercat/redis"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Load(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    redis := redis2.GetClient()
    if redis == nil {
        return
    }
    defer redis.Close()

    var pic string
    if pic, err = redis.SRandMember(KEY).Result(); err != nil {
        logger("LOAD SMEMBERS", err.Error())
    }

    msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, pic)
    if _, err = bot.Send(msg); err != nil {
        logger("LOAD SEND", err.Error())
    }
}
