package handlerDonate

import (
    "embercat/bot/types"
    redis2 "embercat/redis"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

func Show(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    redis := redis2.GetClient()
    if redis == nil {
        return
    }
    defer redis.Close()

    donates := redis.ZRevRangeWithScores(types.REDIS_SUPPORTERS_COLLECTION, 0, -1)
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, getDonatesList(donates))
    msg.ParseMode = tgbotapi.ModeHTML
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Show error %s\n", err.Error())
    }
}
