package handlerPic

import (
    redis2 "embercat/redis"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const KEY = "pic:anime"

func Save(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatID := update.Message.Chat.ID
    reply := update.Message.ReplyToMessage

    if reply == nil || reply.Photo == nil {
        msg := tgbotapi.NewMessage(chatID, "Вызов этой команды возможен только в ответ на сообщение с картинкой")

        if _, err := bot.Send(msg); err != nil {
            logger("bot send", err.Error())
        }
        return
    }

    var photoID string
    for _, pic := range *reply.Photo {
        photoID = pic.FileID
        break
    }

    var err error
    redis := redis2.GetClient()
    if redis == nil {
        return
    }
    defer redis.Close()

    if _, err = redis.SAdd(KEY, photoID).Result(); err != nil {
        logger("REDIS SADD", err.Error())
        return
    }
}
