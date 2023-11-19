package handlerDonate

import (
    "embercat/bot/types"
    redis2 "embercat/redis"
    "log"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Add(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    userID := update.Message.From.UserName
    if userID != "tatoshko" {
        log.Print("Add error this user cannot donate\n")
        return
    }

    redis := redis2.GetClient()
    if redis == nil {
        return
    }
    defer redis.Close()

    info, err := newDonateInfo(update.Message.CommandArguments())
    if err != nil {
        log.Printf("Add error %s\n", err.Error())
        return
    }

    redis.ZAdd(types.REDIS_SUPPORTERS_COLLECTION, info)
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, getDonateMessage(info))
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Add error %s\n", err.Error())
    }
}
