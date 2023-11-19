package handlerPic

import (
    redis2 "embercat/redis"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "strings"
)

func Load(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    redis := redis2.GetClient()
    if redis == nil {
        return
    }
    defer redis.Close()

    var res []string
    if res, err = redis.SMembers(KEY).Result(); err != nil {
        logger("LOAD SMEMBERS", err.Error())
    }

    msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(res, "\n"))
    msg.ParseMode = tgbotapi.ModeHTML

    if _, err = bot.Send(msg); err != nil {
        logger("LOAD SEND", err.Error())
    }
}
